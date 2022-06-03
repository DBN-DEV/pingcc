package app

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/hhyhhy/tsdb"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"go.uber.org/zap"

	"pingcc/domain"
	"pingcc/log"
)

const _aggInterval = 30 * time.Second

type Aggregator struct {
	agentRepo domain.AgentRepo
	tsdb      *tsdb.TSDB[PingResult]
	influxdb  api.WriteAPI
	logger    *zap.Logger
}

func NewAggregator(agentRepo domain.AgentRepo, tsdb *tsdb.TSDB[PingResult], writeAPI api.WriteAPI) *Aggregator {
	logger := log.LWithSvcName("Aggregator")
	a := &Aggregator{agentRepo: agentRepo, tsdb: tsdb, influxdb: writeAPI, logger: logger}

	return a
}

func (a *Aggregator) AggProc() {
	ticker := time.NewTicker(_aggInterval)
	for {
		select {
		case <-ticker.C:
			go a.aggPingResult()
		}
	}
}

func (a *Aggregator) aggPingResult() {
	agents, err := a.agentRepo.AllWithPingTargets(context.Background())
	if err != nil {
		a.logger.Info("Fail to get agents", zap.Error(err))
		return
	}

	for _, agent := range agents {
		a.aggSingeAgentPingResult(agent)
	}
}

func (a *Aggregator) aggSingeAgentPingResult(agent domain.Agent) {
	// 找到每个 agent 的去除ip后的所有的系列，针对每个系列进行聚合
	series := make(map[string][]tsdb.Tag)
	for _, target := range agent.PingTargets {
		s := target.SeriesWithoutIP()
		if _, ok := series[s]; ok {
			continue
		}

		series[s] = target.TagsWithoutIP()
	}

	max := time.Now()
	min := max.Add(-_aggInterval)
	// 计算每个去除了ip的系列的丢包率和rtt
	for _, tags := range series {
		tags = append(tags, tsdb.Tag{Key: "measurement", Value: "icmp"})
		values := a.tsdb.QueryPoints(tags, min, max)
		// tsdb 会将带有 ip 的系列返回，这里进行聚合计算
		avgRtt, avgLoss := calcAvg(values)

		var influxdbTags map[string]string
		for _, tag := range tags[:len(tags)-1] {
			influxdbTags[tag.Key] = tag.Value
		}
		fields := map[string]interface{}{"avg_rtt": avgRtt, "avg_loss": avgLoss}
		point := influxdb2.NewPoint("icmp", influxdbTags, fields, max)
		a.influxdb.WritePoint(point)
	}
}

func calcAvg(seriesValues map[string][]tsdb.Value[PingResult]) (uint32, float32) {
	var sumRttMicros, count, sumLoss uint32
	for _, values := range seriesValues {
		for _, value := range values {
			sumRttMicros += value.V.RttMicros
			count += 1
			if value.V.IsTimeout {
				sumLoss += 1
			}
		}
	}

	var avgLoss float32
	if sumLoss != 0 {
		avgLoss = float32(sumLoss) / float32(count)
	}
	if sumLoss == count {
		avgLoss = 1
	}

	var avgRttMicros uint32
	if sumLoss != count {
		avgRttMicros = sumRttMicros / (count - sumLoss)
	}

	return avgRttMicros, avgLoss
}

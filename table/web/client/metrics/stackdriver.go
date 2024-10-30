package metrics

import (
	monitoring "cloud.google.com/go/monitoring/apiv3"
	"context"
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"time"
)

var Client *monitoring.MetricClient

func Init() error {
	ctx := context.Background()
	var err error
	Client, err = monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	periodicallyWriteMetrics()
	return nil
}

func periodicallyWriteMetrics() {
	ticker := time.NewTicker(conf.Props.GKE.Metrics.PushDuration)
	go func() {
		for range ticker.C {
			err := writeMetric("broadcast/connection-count", 0)
			if err != nil {
				log.Warnf("Couldn't send metric to gke monitoring: %s", err)
			}
		}
	}()
}

func writeMetric(name string, value int64) error {
	now := &timestamp.Timestamp{
		Seconds: time.Now().Unix(),
	}
	req := &monitoringpb.CreateTimeSeriesRequest{
		Name: "projects/" + conf.Props.GKE.ProjectID,
		TimeSeries: []*monitoringpb.TimeSeries{{
			Metric: &metricpb.Metric{
				Type: buildMetricType(name),
			},
			Resource: &monitoredres.MonitoredResource{
				Type: "k8s_pod",
				Labels: map[string]string{
					"project_id":     conf.Props.GKE.ProjectID,
					"location":       conf.Props.GKE.Location,
					"cluster_name":   conf.Props.GKE.ClusterName,
					"namespace_name": conf.Props.GKE.NamespaceName,
					"pod_name":       conf.Props.GKE.PodName,
				},
			},
			Points: []*monitoringpb.Point{{
				Interval: &monitoringpb.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_Int64Value{
						Int64Value: value,
					},
				},
			}},
		}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := Client.CreateTimeSeries(ctx, req)
	if err != nil {
		return fmt.Errorf("could not write time series value, %v ", err)
	}
	return nil
}

func buildMetricType(name string) string {
	return fmt.Sprintf("custom.googleapis.com/table/%s", name)
}

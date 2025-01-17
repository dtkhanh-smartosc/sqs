package domain

import "github.com/prometheus/client_golang/prometheus"

var (
	// sqs_ingest_usecase_process_block_duration
	//
	// gauge that measures the duration of processing a block in milliseconds in ingest usecase
	//
	// Has the following labels:
	// * height - the height of the block being processed
	SQSIngestUsecaseProcessBlockDurationMetricName = "sqs_ingest_usecase_process_block_duration"

	// sqs_ingest_usecase_process_block_error
	//
	// counter that measures the number of errors that occur during processing a block in ingest usecase
	//
	// Has the following labels:
	// * err - the error message occurred
	// * height - the height of the block being processed
	SQSIngestUsecaseProcessBlockErrorMetricName = "sqs_ingest_usecase_process_block_error_total"

	// sqs_ingest_usecase_parse_pool_error_total
	//
	// counter that measures the number of errors that occur during pool parsing in ingest usecase
	//
	// Has the following labels:
	// * err - the error message occurred
	SQSIngestUsecaseParsePoolErrorMetricName = "sqs_ingest_usecase_parse_pool_error_total"

	// sqs_pricing_worker_compute_error_counter
	//
	// counter that measures the number of errors that occur during pricing worker computation
	//
	// Has the following labels:
	// * height - the height of the block being processed
	SQSPricingWorkerComputeErrorCounterMetricName = "sqs_pricing_worker_compute_error_total"

	// sqs_pricing_worker_compute_duration
	//
	// gauge that tracks duration of pricing worker computation
	SQSPricingWorkerComputeDurationMetricName = "sqs_pricing_worker_compute_duration"

	// sqs_pool_liq_pricing_worker_compute_duration
	//
	// gauge that tracks duration of pricing worker computation
	SQSPoolLiquidityPricingWorkerComputeDurationMetricName = "sqs_pool_liq_pricing_worker_compute_duration"

	// sqs_update_assets_at_block_height_interval_error_total
	//
	// counter that measures the number of errors that occur during updating assets in ingest usecase
	// Has the following labels:
	// * height - the height of the block being processed
	SQSUpdateAssetsAtHeightIntervalMetricName = "sqs_update_assets_at_block_height_interval_error_total"

	// sqs_pricing_errors_total
	//
	// counter that measures the number of pricing errors
	// Has the following labels:
	// * base - the base asset symbol
	// * quote - the quote asset symbol
	// * err - the error message occurred
	SQSPricingErrorCounterMetricName = "sqs_pricing_errors_total"

	// sqs_pricing_fallback_total
	//
	// counter that measures the number of fallback from chain pricing source to coingecko
	// Has the following labels:
	// * base - the base asset symbol
	// * quote - the quote asset symbol
	SQSPricingFallbackCounterMetricName = "sqs_pricing_fallback_total"

	// sqs_passthrough_numia_aprs_fetch_error_total
	//
	// counter that measures the number of errors when fetching APRs from Numia in a passthrough module.
	SQSPassthroughNumiaAPRsFetchErrorCounterMetricName = "sqs_passthrough_numia_aprs_fetch_error_total"

	// sqs_passthrough_timeseries_pool_fees_fetch_error_total
	//
	// counter that measures the number of errors when fetching fees from timeseries data stack in a passthrough module.
	SQSPassthroughTimeseriesPoolFeesFetchErrorCounterMetricName = "sqs_passthrough_timeseries_pool_fees_fetch_error_total"

	SQSIngestHandlerProcessBlockDurationGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: SQSIngestUsecaseProcessBlockDurationMetricName,
			Help: "gauge that measures the duration of processing a block in milliseconds in ingest usecase",
		},
	)

	SQSIngestHandlerProcessBlockErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: SQSIngestUsecaseProcessBlockErrorMetricName,
			Help: "counter that measures the number of errors that occur during processing a block in ingest usecase",
		},
		[]string{"err", "height"},
	)

	SQSIngestHandlerPoolParseErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: SQSIngestUsecaseParsePoolErrorMetricName,
			Help: "counter that measures the number of errors that occur during pool parsing in ingest usecase",
		},
		[]string{"err"},
	)

	SQSPricingWorkerComputeErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: SQSPricingWorkerComputeErrorCounterMetricName,
			Help: "counter that measures the number of errors that occur during pricing worker computation",
		},
		[]string{"height"},
	)

	SQSPricingWorkerComputeDurationGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: SQSPricingWorkerComputeDurationMetricName,
			Help: "gauge that tracks duration of pricing worker computation",
		},
	)

	SQSPoolLiquidityPricingWorkerComputeDurationGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: SQSPoolLiquidityPricingWorkerComputeDurationMetricName,
			Help: "gauge that tracks duration of pool liquidity pricing worker computation",
		},
	)

	SQSUpdateAssetsAtHeightIntervalErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: SQSUpdateAssetsAtHeightIntervalMetricName,
			Help: "Update assets at block height interval error when processing block data",
		},
		[]string{"err", "height"},
	)

	SQSPricingErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: SQSPricingErrorCounterMetricName,
			Help: "Total number of pricing errors",
		},
		[]string{"base", "quote", "err"},
	)
	SQSPricingFallbackCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: SQSPricingFallbackCounterMetricName,
			Help: "Total number of fallback from chain pricing source to coingecko",
		},
		[]string{"base", "quote"},
	)

	SQSPassthroughNumiaAPRsFetchErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: SQSPassthroughNumiaAPRsFetchErrorCounterMetricName,
			Help: "Total number of errors when fetching APRs from Numia in a passthrough module.",
		},
	)

	SQSPassthroughTimeseriesPoolFeesFetchErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: SQSPassthroughTimeseriesPoolFeesFetchErrorCounterMetricName,
			Help: "Total number of errors when fetching pool fees from timeseries in a passthrough module.",
		},
	)
)

func init() {
	prometheus.MustRegister(SQSIngestHandlerProcessBlockDurationGauge)
	prometheus.MustRegister(SQSIngestHandlerProcessBlockErrorCounter)
	prometheus.MustRegister(SQSIngestHandlerPoolParseErrorCounter)
	prometheus.MustRegister(SQSPricingWorkerComputeDurationGauge)
	prometheus.MustRegister(SQSPricingWorkerComputeErrorCounter)
	prometheus.MustRegister(SQSPoolLiquidityPricingWorkerComputeDurationGauge)
	prometheus.MustRegister(SQSUpdateAssetsAtHeightIntervalErrorCounter)
	prometheus.MustRegister(SQSPricingErrorCounter)
	prometheus.MustRegister(SQSPricingFallbackCounter)
	prometheus.MustRegister(SQSPassthroughNumiaAPRsFetchErrorCounter)
	prometheus.MustRegister(SQSPassthroughTimeseriesPoolFeesFetchErrorCounter)
}

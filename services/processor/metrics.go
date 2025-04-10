package main

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Métricas customizadas
var (
    logsProcessedTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "logs_processed_total",
            Help: "Total de logs processados com sucesso",
        },
        []string{"level", "service"},
    )

    logsErroredTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "logs_errors_total",
            Help: "Total de erros ao processar logs",
        },
    )
)

func initMetrics() {
    prometheus.MustRegister(logsProcessedTotal)
    prometheus.MustRegister(logsErroredTotal)

    // Servidor HTTP pra expor as métricas
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
    }()
}

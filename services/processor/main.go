package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/google/uuid"
    "github.com/segmentio/kafka-go"
)

// Estrutura recebida do Kafka
type IncomingLog struct {
    Service   string `json:"service"`
    Message   string `json:"message"`
    Level     string `json:"level"`
    Timestamp string `json:"timestamp,omitempty"`
}

// Estrutura enriquecida que será enviada para o Elasticsearch
type EnrichedLog struct {
    ID        string `json:"id"`
    Service   string `json:"service"`
    Message   string `json:"message"`
    Level     string `json:"level"`
    Timestamp string `json:"timestamp"`
    Hostname  string `json:"hostname"`
}

func enrich(log IncomingLog) EnrichedLog {
    id := uuid.New().String()
    hostname, _ := os.Hostname()

    ts := log.Timestamp
    if ts == "" {
        ts = time.Now().UTC().Format(time.RFC3339)
    }

    return EnrichedLog{
        ID:        id,
        Service:   log.Service,
        Message:   log.Message,
        Level:     log.Level,
        Timestamp: ts,
        Hostname:  hostname,
    }
}

func sendToElasticsearch(logData EnrichedLog) {
    esURL := os.Getenv("ELASTICSEARCH_URL")
    if esURL == "" {
        esURL = "http://elasticsearch:9200"
    }

    jsonBody, err := json.Marshal(logData)
    if err != nil {
        logsErroredTotal.Inc()
        log.Printf("❌ Erro ao serializar log: %v\n", err)
        return
    }

    req, err := http.NewRequest("POST", fmt.Sprintf("%s/logs/_doc", esURL), bytes.NewBuffer(jsonBody))
    if err != nil {
        logsErroredTotal.Inc()
        log.Printf("❌ Erro ao criar request: %v\n", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        logsErroredTotal.Inc()
        log.Printf("❌ Erro ao enviar para Elasticsearch: %v\n", err)
        return
    }
    defer resp.Body.Close()

    log.Printf("📤 Enviado para Elasticsearch (status: %s)\n", resp.Status)
}

func main() {
    kafkaBroker := os.Getenv("KAFKA_BROKER")
    if kafkaBroker == "" {
        kafkaBroker = "localhost:9092"
    }

    initMetrics() // Inicializa métricas Prometheus

    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{kafkaBroker},
        Topic:   "logs",
        GroupID: "log-processor-group",
    })
    defer reader.Close()

    log.Println("🚀 Processor pronto para consumir do tópico 'logs'")

    for {
        m, err := reader.ReadMessage(context.Background())
        if err != nil {
            logsErroredTotal.Inc()
            log.Printf("❌ Erro ao ler mensagem: %v\n", err)
            continue
        }

        var incoming IncomingLog
        if err := json.Unmarshal(m.Value, &incoming); err != nil {
            logsErroredTotal.Inc()
            log.Printf("❌ Erro ao parsear JSON: %v\n", err)
            continue
        }

        enriched := enrich(incoming)
        logsProcessedTotal.WithLabelValues(enriched.Level, enriched.Service).Inc()

        sendToElasticsearch(enriched)
    }
}

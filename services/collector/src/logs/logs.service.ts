import { Injectable, OnModuleInit, Logger } from '@nestjs/common';
import { CreateLogDto } from './dto/create-log.dto';
import { Kafka, Producer } from 'kafkajs';

@Injectable()
export class LogsService implements OnModuleInit {
  private kafka = new Kafka({ brokers: [process.env.KAFKA_BROKER || 'localhost:9092'] });
  private producer!: Producer;
  private readonly logger = new Logger(LogsService.name);

  async onModuleInit() {
    this.producer = this.kafka.producer();
    await this.producer.connect();
  }

  async publish(log: CreateLogDto) {
    await this.producer.send({
      topic: 'logs',
      messages: [{ value: JSON.stringify(log) }],
    });

    this.logger.log(`Publicado no Kafka: ${JSON.stringify(log)}`);
  }
}

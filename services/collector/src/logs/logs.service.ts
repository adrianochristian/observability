import { Injectable, OnModuleInit } from '@nestjs/common';
import { CreateLogDto } from './dto/create-log.dto';
import { Kafka, Producer } from 'kafkajs';

@Injectable()
export class LogsService implements OnModuleInit {
  private kafka = new Kafka({ brokers: [process.env.KAFKA_BROKER || 'localhost:9092'] });
  private producer!: Producer;

  async onModuleInit() {
    this.producer = this.kafka.producer();
    await this.producer.connect();
  }

  async publish(log: CreateLogDto) {
    await this.producer.send({
      topic: 'logs',
      messages: [{ value: JSON.stringify(log) }],
    });
    console.log('Publicado no Kafka:', log);
  }
}

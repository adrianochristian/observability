import { Injectable } from '@nestjs/common';
import { CreateLogDto } from './dto/create-log.dto';

@Injectable()
export class LogsService {
  async publish(log: CreateLogDto) {
    console.log('Log recebido:', log);
    // TODO: publicar no Kafka
  }
}

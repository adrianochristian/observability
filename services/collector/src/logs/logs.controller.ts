import { Controller, Post, Body } from '@nestjs/common';
import { LogsService } from './logs.service';
import { CreateLogDto } from './dto/create-log.dto';

@Controller('logs')
export class LogsController {
  constructor(private readonly logsService: LogsService) {}

  @Post()
  async receiveLog(@Body() log: CreateLogDto): Promise<{ status: string }> {
    await this.logsService.publish(log);
    return { status: 'received' };
  }
}

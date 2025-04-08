import { IsString, IsIn, IsOptional } from 'class-validator';

export class CreateLogDto {
  @IsString()
  service!: string;

  @IsString()
  message!: string;

  @IsIn(['debug', 'info', 'warn', 'error', 'fatal'])
  level!: string;

  @IsOptional()
  @IsString()
  timestamp?: string;
}

import { Injectable } from '@nestjs/common';
import * as bcrypt from 'bcrypt';

export type User = {
  id: number;
  username: string;
  password: string;
  role: 'admin' | 'viewer';
};

@Injectable()
export class UsersService {
  private readonly users: User[] = [
    {
      id: 1,
      username: 'admin',
      password: bcrypt.hashSync('admin123', 10),
      role: 'admin',
    },
    {
      id: 2,
      username: 'viewer',
      password: bcrypt.hashSync('viewer123', 10),
      role: 'viewer',
    },
  ];

  async findByUsername(username: string): Promise<User | undefined> {
    return this.users.find(u => u.username === username);
  }
}
import { Role } from '../enums/roles.enum';

export interface UserAuth {
  id: string;
  username: string;
  email: string;
  role: Role;
  confirmed: boolean;
}

import { Role } from "../enums/roles.enum";

export interface User {
    id: string,
    firstName: string,
    lastName: string,
    email: string,
    residence: string,
    role: Role,
    username: string,
    age: number,
}
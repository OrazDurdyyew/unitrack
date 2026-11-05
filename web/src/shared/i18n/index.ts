import ru from './ru.json';

export const messages = { ru } as const;
export type MessageSchema = typeof ru;
export const defaultLocale = 'ru';

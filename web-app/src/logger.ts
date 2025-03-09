import pino, { Logger } from "pino"; 
import 'dotenv/config';

let log: Logger;

if (process.env.PRETTY_LOGGER) {
    log = pino({
        transport: {
            target: 'pino-pretty',
            options: {
                colorize: true,
                singleLine: true,
            },
        },
    });
} else {
    log = pino({});
}

export default log;

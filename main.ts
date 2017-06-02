'use strict';

import * as Hapi from 'hapi';
import * as Joi from 'joi';

const server: Hapi.Server = new Hapi.Server()
server.connection({ port: 3000 });

server.route({
    method: 'GET',
    path: '/',
    handler: (request: Hapi.Request, reply: Hapi.ReplyNoContinue) => {
        reply('Hello World\n')
    }
});

server.route({
    method: 'GET',
    path: '/hello/{name}',
    handler: (request: Hapi.Request, reply: Hapi.ReplyNoContinue) => {
        reply(`Hello, ${encodeURIComponent(request.params.name)}!\n`);
    },
    config: {
        validate: {
            params: {
                name: Joi.string().min(3).max(10)
            }
        }
    }
});

server.start((err) => {
    if (err) {
        throw err;
    }
    console.log('server running at 3000');
})

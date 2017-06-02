'use strict';
exports.__esModule = true;
var Hapi = require("hapi");
var Joi = require("joi");
var server = new Hapi.Server();
server.connection({ port: 3000 });
server.route({
    method: 'GET',
    path: '/',
    handler: function (request, reply) {
        reply('Hello World\n');
    }
});
server.route({
    method: 'GET',
    path: '/hello/{name}',
    handler: function (request, reply) {
        reply("Hello, " + encodeURIComponent(request.params.name) + "!\n");
    },
    config: {
        validate: {
            params: {
                name: Joi.string().min(3).max(10)
            }
        }
    }
});
server.start(function (err) {
    if (err) {
        throw err;
    }
    console.log('server running at 3000');
});

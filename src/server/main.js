"use strict";
exports.__esModule = true;
var Hapi = require("hapi");
// Create a server with a host and port
var server = new Hapi.Server();
server.connection({
    host: 'localhost',
    port: 8000
});
function comicHandler(request, reply) {
    var body = request.payload;
    return reply(body.actions[0].value);
}
// Add the route
server.route({
    method: 'GET',
    path: '/hello',
    handler: function (request, reply) {
        return reply('hello world');
    }
});
server.route({
    method: 'POST',
    path: '/operations',
    handler: function (request, reply) {
        var body = request.payload;
        var route = body.callback_id.split('_')[0];
        if (route === 'comic') {
            return comicHandler(request, reply);
        }
    }
});
server.route({
    method: 'GET',
    path: '/comic',
    handler: comicHandler
});
// Start the server
server.start(function (err) {
    if (err) {
        throw err;
    }
    console.log('Server running at:', server.info.uri);
});

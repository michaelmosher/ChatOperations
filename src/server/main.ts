import * as Hapi from 'hapi'

// Create a server with a host and port
const server = new Hapi.Server();
server.connection({ 
    host: 'localhost', 
    port: 8000 
});

interface SlackAction{
    name: string,
    value: string,
    type: string,
}

interface SlackInteractiveResponse{
    actions: SlackAction[],
    callback_id: string,
    team: {id: string, domain: string},
    channel: {id: string, name: string},
    user: {id: string, name: string},
    token: string
    response_url: string
}

function comicHandler (request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
    let body: SlackInteractiveResponse = request.payload
    return reply(body.actions[0].value)
}

// Add the route
server.route({
    method: 'GET',
    path:'/hello', 
    handler: function (request, reply) {

        return reply('hello world');
    }
});

server.route({
    method: 'POST',
    path: '/operations',
    handler: function (request, reply) {
        let body: SlackInteractiveResponse = JSON.parse(request.payload.payload)
        let route = body.callback_id.split('_')[0]
        if (route === 'comic') {
            return comicHandler(request, reply)
        }
    }
})

server.route({
    method: 'GET',
    path: '/comic',
    handler: comicHandler
})

// Start the server
server.start((err) => {

    if (err) {
        throw err;
    }
});
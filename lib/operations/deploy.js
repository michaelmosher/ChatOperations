"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var Server = (function () {
    function Server(name, address) {
        this.name = name;
        this.address = address;
    }
    return Server;
}());
exports.Server = Server;
var Request = (function () {
    function Request(requester) {
        this.requester = requester;
        var suffix = String(Date.now() % 100000000); // random enough
        this.action = 'deploy';
        this.callback_id = "deploy_" + requester + "_" + suffix;
    }
    Request.prototype.update = function (u) {
        switch (u.constructor.name) {
            case 'Server': {
                this.server = u;
                break;
            }
            case 'Responder': {
                this.responder = u.value;
                break;
            }
            case 'Approval': {
                this.approved = u.value;
                break;
            }
            case 'Success': {
                this.succeeded = u.value;
                break;
            }
            case 'ResponseURL': {
                this.response_url = u.value;
                break;
            }
        }
    };
    Request.prototype.isReady = function () {
        return this.server !== undefined
            && this.response_url !== undefined;
    };
    Request.prototype.next = function () {
        // TODO load from repository
        var server = new Server('dev', '123.4.56.789');
        return [server];
    };
    Request.prototype.summary = function () {
        return (this.server != undefined)
            ? this.requester + " has requested a deploy on " + this.server.name + "."
            : this.requester + " has requested a deploy.";
    };
    return Request;
}());
exports.Request = Request;
function init(requester) {
    return new Request(requester);
}
exports.init = init;
//# sourceMappingURL=deploy.js.map
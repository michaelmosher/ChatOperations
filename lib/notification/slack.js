"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var Templates = require("./templates");
var Notifier = (function () {
    function Notifier(http, hostingWebhook) {
        this.http = http;
        this.hostingWebhook = hostingWebhook;
    }
    Notifier.prototype.reportError = function (url, err) {
        var body = Templates.errorReport;
        body.attachments[0].text = err.message;
        return this.http.post({ url: url, body: body });
    };
    Notifier.prototype.handleErrors = function (url) {
        var _this = this;
        return function (err, req, body) {
            if (err) {
                return _this.reportError(url, err);
            }
            if (req.statusCode !== 200) {
                return _this.reportError(url, body.error);
            }
        };
    };
    Notifier.prototype.requestSubmitted = function (r) {
        var body = Templates.submitRequest;
        body.text = r.summary();
        body.attachments[0].callback_id = r.id;
        this.http.post({
            url: this.hostingWebhook,
            body: body
        }, this.handleErrors(r.response_url));
    };
    Notifier.prototype.requestAnswered = function (r, approved) {
        var body = {
            text: '',
            replace_original: false
        };
        body.text = (approved)
            ? ":white_check_mark: " + r.responder + " approved your request."
            : ":x: " + r.responder + " denied your request.";
        this.http.post({
            url: r.response_url,
            body: body
        }, this.handleErrors(this.hostingWebhook));
    };
    Notifier.prototype.requestApproved = function (r) {
        this.requestAnswered(r, true);
    };
    Notifier.prototype.requestDenied = function (r) {
        this.requestAnswered(r, false);
    };
    return Notifier;
}());
exports.Notifier = Notifier;
//# sourceMappingURL=slack.js.map
"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = y[op[0] & 2 ? "return" : op[0] ? "throw" : "next"]) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [0, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var Templates = require("./templates");
var Notifier = (function () {
    function Notifier(http, hostingWebhook) {
        this.http = http;
        this.hostingWebhook = hostingWebhook;
    }
    Notifier.prototype.promPostRequest = function (options) {
        return __awaiter(this, void 0, void 0, function () {
            var _this = this;
            return __generator(this, function (_a) {
                return [2 /*return*/, new Promise(function (resolve, reject) {
                        _this.http.post(options, function (err, req, body) {
                            if (err) {
                                return reject(err);
                            }
                            if (req.statusCode !== 200) {
                                return reject(body.error);
                            }
                            return resolve(body);
                        });
                    })];
            });
        });
    };
    Notifier.prototype.reportError = function (url, err) {
        var body = Templates.errorReport;
        body.attachments[0].text = err.message;
        return this.promPostRequest({ url: url, body: body });
    };
    Notifier.prototype.requestSubmitted = function (r) {
        var _this = this;
        var body = Templates.submitRequest;
        body.text = r.summary();
        body.attachments[0].callback_id = r.callback_id;
        return this.promPostRequest({
            url: this.hostingWebhook,
            body: body
        })
            .catch(function (error) {
            _this.reportError(r.response_url, error);
        });
    };
    Notifier.prototype.requestAnswered = function (r, approved) {
        var _this = this;
        var body = {
            text: '',
            replace_original: false
        };
        body.text = (approved)
            ? ":white_check_mark: " + r.responder + " approved your request."
            : ":x: " + r.responder + " denied your request.";
        return this.promPostRequest({
            url: r.response_url,
            body: body
        })
            .catch(function (error) {
            _this.reportError(_this.hostingWebhook, error);
        });
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
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.errorReport = {
    text: "Something seems to have gone wrong. Please report this to @hostingteam.",
    attachments: [{
            text: "PLACEHOLDER",
            color: "red"
        }]
};
exports.submitRequest = {
    text: "PLACEHOLDER",
    attachments: [
        {
            text: "Please approve or reject the request:",
            fallback: "Please use an official Slack client for Ops help",
            callback_id: 0,
            attachment_type: "default",
            actions: [
                {
                    name: "ops_request_submitted",
                    text: "Approve",
                    style: "primary",
                    type: "button",
                    value: "approved"
                },
                {
                    name: "ops_request_submitted",
                    text: "Deny",
                    style: "danger",
                    type: "button",
                    value: "rejected"
                }
            ]
        }
    ]
};
//# sourceMappingURL=templates.js.map
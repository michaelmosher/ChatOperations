export class Responder {
	constructor(public value: string) {}
}

export class Approval {
	constructor(public value: boolean) {}
}

export class Success {
	constructor(public value: boolean) {}
}

export class ResponseURL {
	constructor(public value: string) {}
}

export interface Request {
	callback_id:  string
	action:       string
	requester:    string
	responder:    string
	approved:     boolean
	succeeded:    boolean
	response_url: string
	update(u : any): void
	isReady(): boolean
	next(): Array<any>
	summary(): string
}

export interface RequestRepo {
	findById(callback_id: string): Promise<Request>
	store(r: Request): Promise<Request>
}

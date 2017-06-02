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
	id:           number
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
	findAllTypes(): Promise<Array<string>>
	findById(id: number): Promise<Request>
	store(r: Request): Promise<Request>
}

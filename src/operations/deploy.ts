export class Server {
	constructor(readonly name: string, readonly address: string) {}
}

export class Request {
	[k: string]: any;
	readonly action: string
	callback_id:  string
	responder:    string
	approved:     boolean
	succeeded:    boolean
	response_url: string
	server:       Server

    constructor(public requester: string) {
		let suffix = String(Date.now() % 100000000) // random enough
		this.action = 'deploy'
		this.callback_id = `deploy_${requester}_${suffix}`
    }

	update(u: any) : void {
		switch (u.constructor.name) {
			case 'Server': {
				this.server = <Server>u
				break
			}
			case 'Responder': {
				this.responder = u.value
				break
			}
			case 'Approval': {
				this.approved = u.value
				break
			}
			case 'Success': {
				this.succeeded = u.value
				break
			}
			case 'ResponseURL': {
				this.response_url = u.value
				break
			}
		}
	}

	isReady() : boolean {
		return this.server !== undefined
			&& this.response_url !== undefined
	}

	next() : Array<any> {
		// TODO load from repository
		let server = new Server('dev', '123.4.56.789')
		return [server]
	}

	summary() : string {
		return (this.server != undefined)
		  ? `${this.requester} has requested a deploy on ${this.server.name}.`
		  : `${this.requester} has requested a deploy.`
	}
}

export function init(requester : string) : Request {
	return new Request(requester)
}

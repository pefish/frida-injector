export namespace main {
	
	export class ProcessInfo {
	    name: string;
	    id: number;
	    parent_id: number;
	    user: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.id = source["id"];
	        this.parent_id = source["parent_id"];
	        this.user = source["user"];
	        this.path = source["path"];
	    }
	}

}


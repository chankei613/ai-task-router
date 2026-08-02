export namespace api {
	
	export class IssueKeyResult {
	    id: string;
	    name: string;
	    api_key: string;
	
	    static createFrom(source: any = {}) {
	        return new IssueKeyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.api_key = source["api_key"];
	    }
	}
	export class RouteResult {
	    log: db.RoutingLog;
	    chosen_model?: db.ModelSpec;
	
	    static createFrom(source: any = {}) {
	        return new RouteResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.log = this.convertValues(source["log"], db.RoutingLog);
	        this.chosen_model = this.convertValues(source["chosen_model"], db.ModelSpec);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace db {
	
	export class AgentKey {
	    id: string;
	    name: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    revoked_at?: any;
	
	    static createFrom(source: any = {}) {
	        return new AgentKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.revoked_at = this.convertValues(source["revoked_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ModelSpec {
	    id: string;
	    provider: string;
	    model_id: string;
	    quality_tier: string;
	    input_price_per_1m: number;
	    output_price_per_1m: number;
	    capabilities: string[];
	    enabled: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ModelSpec(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.provider = source["provider"];
	        this.model_id = source["model_id"];
	        this.quality_tier = source["quality_tier"];
	        this.input_price_per_1m = source["input_price_per_1m"];
	        this.output_price_per_1m = source["output_price_per_1m"];
	        this.capabilities = source["capabilities"];
	        this.enabled = source["enabled"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RoutingLog {
	    id: string;
	    source: string;
	    task_type: string;
	    required_capabilities: string[];
	    min_quality_tier: string;
	    chosen_model_id: string;
	    reasoning: string;
	    // Go type: time
	    requested_at: any;
	
	    static createFrom(source: any = {}) {
	        return new RoutingLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.source = source["source"];
	        this.task_type = source["task_type"];
	        this.required_capabilities = source["required_capabilities"];
	        this.min_quality_tier = source["min_quality_tier"];
	        this.chosen_model_id = source["chosen_model_id"];
	        this.reasoning = source["reasoning"];
	        this.requested_at = this.convertValues(source["requested_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}


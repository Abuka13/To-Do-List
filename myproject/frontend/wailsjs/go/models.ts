export namespace domain {
	
	export class Task {
	    id: number;
	    title: string;
	    isCompleted: boolean;
	    createdAt: string;
	    completedAt?: string;
	    dueDate?: string;
	    priority: string;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.isCompleted = source["isCompleted"];
	        this.createdAt = source["createdAt"];
	        this.completedAt = source["completedAt"];
	        this.dueDate = source["dueDate"];
	        this.priority = source["priority"];
	    }
	}

}


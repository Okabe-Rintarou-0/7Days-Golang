export function post(url: string, json: Object, callback: Function, onError?: Function) {
    let opts: RequestInit = {
        method: "POST",
        body: JSON.stringify(json),
        headers: {
            'Content-Type': 'application/json'
        },
    };

    fetch(url, opts)
        .then(response => response.json())
        .then(data => {
            callback(data);
        })
        .catch(error => {
            if (onError !== undefined)
                onError(error)
        })
}

export function put(url: string, callback: Function, onError?: Function) {
    let opts: RequestInit = {
        method: "PUT",
        headers: {
            'Content-Type': 'application/json'
        },
    };
    fetch(url, opts)
        .then(response => response.text())
        .then(data => {
            callback(data);
        })
        .catch(error => {
            if (onError !== undefined)
                onError(error)
        })
}

export function del(url: string, callback: Function, onError?: Function) {
    let opts: RequestInit = {
        method: "DELETE",
    };
    fetch(url, opts)
        .then(response => response.json())
        .then(data => {
            callback(data);
        })
        .catch(error => {
            if (onError !== undefined)
                onError(error)
        })
}


export function get(url: string, callback: Function, onError?: Function) {
    let opts: RequestInit = {
        method: "GET",
    };
    fetch(url, opts)
        .then(response => response.json())
        .then(data => {
            callback(data);
        })
        .catch(error => {
            if (onError !== undefined)
                onError(error)
        })
}

export function getForText(url: string, callback: Function, onError?: Function) {
    let opts: RequestInit = {
        method: "GET",
    };
    fetch(url, opts)
        .then(response => response.text())
        .then(response => callback(response))
        .catch(error => {
            if (onError !== undefined)
                onError(error)
        })
}

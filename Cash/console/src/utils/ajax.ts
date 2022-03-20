export function post(url: string, json: Object, callback: Function) {
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
            console.log(error);
        })
}

export function put(url: string, json: Object, callback: Function) {
    let opts: RequestInit = {
        method: "PUT",
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
            console.log(error);
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
            if (onError != undefined)
                onError(error)
        })
}

let dataTimeout = 55000;
let settingsTimeout = 5000;
let iteration = 0;
let countBlock = top.document.getElementById('count');
let errorBlock = top.document.getElementById("error");
console.log(countBlock, errorBlock);

async function getSettings() {
    const url = "/settings.json";
    try {
        const response = await fetch(url);
        if (!response.ok) {
            console.log("!ok response settings", response);
            errorBlock.style.display = 'block';
            errorBlock.textContent = "Error settings connect " + response
        }
        const json = await response.json();
        if (json === undefined || json === null) {
            console.log("json settings undefined");
        }
        if (json.reboot === true && iteration !== 1) {
            window.location.reload();
        }
        if (json.timeout !== undefined && json.timeout != 0) {
            dataTimeout = json.timeout
        }
    } catch (error) {
        console.error(error.message);
        errorBlock.style.display = 'block';
        errorBlock.textContent = "Error api connect " + error.message
    }
    setTimeout(getSettings, settingsTimeout);
}

async function getData() {
    iteration++;

    if (countBlock != null) {
        console.log("countBlock null");
        countBlock.textContent = String(iteration)
    }
    const url = "/api.json";
    try {
        const response = await fetch(url);
        if (!response.ok) {
            console.log("!ok response", response);
            errorBlock.style.display = 'block';
            errorBlock.textContent = "Error api connect " + response
        }
        const json = await response.json();
        if (json === undefined || json === null) {
            console.log("json undefined");
        }
        if (json.reboot === true && iteration !== 1) {
            window.location.reload();
        }
        if (json.info !== undefined || json !== null) {
            for (const i in json.info) {
                errorBlock.style.display = 'none';
                console.log(json.info[i]);
                let block = top.document.getElementById(json.info[i].block);

                let time, name, content, error;
                let div = top.document.getElementById(json.info[i].id)
                if (div === undefined || div === null) {
                    div = document.createElement("div");
                    div.setAttribute("class", "item");
                    div.setAttribute("id", json.info[i].id);
                    block.append(div);

                    time = document.createElement("span");
                    name = document.createElement("span");
                    content = document.createElement("span");
                    error = document.createElement("span");

                    time.setAttribute("class", "item-time");
                    name.setAttribute("class", "item-name");
                    content.setAttribute("class", "item-content");
                    error.setAttribute("class", "item-error");

                    time.setAttribute("id", json.info[i].id+"-item-time");
                    name.setAttribute("id", json.info[i].id+"-item-name");
                    content.setAttribute("id", json.info[i].id+"-item-content");
                    error.setAttribute("id", json.info[i].id+"-item-error");

                    div.append(time);
                    div.append(name);
                    div.append(content);
                    div.append(error);
                } else {
                    time = top.document.getElementById(json.info[i].id+"-item-time");
                    name = top.document.getElementById(json.info[i].id+"-item-name");
                    content = top.document.getElementById(json.info[i].id+"-item-content");
                    error = top.document.getElementById(json.info[i].id+"-item-error");
                }
                //console.log("test", div, time);

                time.textContent = json.info[i].time;
                name.textContent = json.info[i].name;
                content.textContent = json.info[i].value;
                error.textContent = json.info[i].error;
            }
        }
        //console.log(json);
    } catch (error) {
        console.error(error.message);
        errorBlock.style.display = 'block';
        errorBlock.textContent = "Error api connect " + error.message
    }
    setTimeout(getData, dataTimeout);
}
getSettings();
getData();

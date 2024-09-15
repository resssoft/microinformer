
let built = false
let dataTimeout = 55000;
let settingsTimeout = 5000;
let iteration = 0;
let countBlock = top.document.getElementById('count');
let errorBlock = top.document.getElementById("error");
console.log(countBlock, errorBlock);
const urlSettings = "/api/settings/data.json";
const urlData = "/api/items/data.json";

async function getSettings() {
    try {
        const response = await fetch(urlSettings);
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
            setTimeout(function () {
                window.location.reload();
            }, 500);
            return
        }
        if (json.timeout !== undefined && json.timeout != 0) {
            dataTimeout = json.timeout
        }
        if (!built) {
            console.log("start build");
            if (json.panel !== undefined && json.panel != 0) {
                let container = top.document.getElementById("content");
                console.log("start rows");
                for (const i in json.panel.rows) {
                    row = document.createElement("div");
                    row.setAttribute("class", "row");
                    row.setAttribute("id", json.panel.rows[i].id);

                    console.log("start blocks");
                    for (const y in json.panel.rows[i].blocks) {
                        block = document.createElement("div");
                        row.setAttribute("class", "block");
                        block.setAttribute("id", json.panel.rows[i].blocks[y].id);
                        row.append(block);
                    }
                    container.append(row);

                }
                built = true;
                getData();
            }
        }
    } catch (error) {
        console.error(error.message);
        errorBlock.style.display = 'block';
        errorBlock.textContent = "Error api connect " + error.message
    }
    setTimeout(getSettings, settingsTimeout);
}

async function getData() {
    console.log("getData");
    if (!built) {
        console.log("skip get data");
        setTimeout(getData, dataTimeout);
        return
    }

    iteration++;

    if (countBlock != null) {
        countBlock.textContent = String(iteration)
    }
    try {
        const response = await fetch(urlData);
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
            setTimeout(function () {
                window.location.reload();
            }, 500);
            return
        }
        if (json.info !== undefined || json !== null) {
            for (const i in json.info) {
                errorBlock.style.display = 'none';
                let block = top.document.getElementById(json.info[i].block);

                let time, name, content, error;
                let div = top.document.getElementById(json.info[i].id);
                if (div === undefined || div === null) {
                    if (json.info[i].once === true) {

                        console.log("show modal");
                        top.document.getElementById("modal").style.display = 'block';
                        top.document.getElementById("modal").textContent = json.info[i].value;
                        setTimeout(hideModal, json.info[i].keep_by * 1000);
                        continue
                    }

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
    } catch (error) {
        console.error(error.message);
        errorBlock.style.display = 'block';
        errorBlock.textContent = "Error api connect " + error.message
    }
    setTimeout(getData, dataTimeout);
}

async function hideModal() {
    console.log("hide modal");
    top.document.getElementById("modal").style.display = 'none';
}
getSettings();

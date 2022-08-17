(function(){
    // let apiPath = 'http://127.0.0.1:8383/test';
    let apiPath = 'http://127.0.0.1:8080';

    let dataScheme = {
        'voice_call' : [
            "country", "bandwidth", "response_time", "provider",
            "connection_stability", "ttfb", "voice_purity",
            "median_of_calls_time"
        ],
        'sms' : ["country", "bandwidth", "response_time", "provider"],
        'mms' : ["country", "bandwidth", "response_time", "provider"],
        'incident' : ["topic", "status"]
    };

    let pieColors = ['#79c779', '#eae537', '#ff4141'];

    let ready = (callback) => {
        if (document.readyState != "loading") callback();
        else document.addEventListener("DOMContentLoaded", callback);
    };

    let showErrors = function(errors) {
        let list = document.createElement('ul');
        errors.forEach((error) => {
            let item = document.createElement('li');
            item.textContent = error;
            list.appendChild(item);
        });
        let container = document.querySelector(".errors-container");
        container.appendChild(list);
        container.style.display = "block";
    };

    let checkJsonScheme = function(json) {
        let keys = ["status", "data", "error"];
        for(i in keys) {
            if(typeof json[keys[i]] === 'undefined') {
                return false;
            }
        }
        if(!json.data && !json.error) {
            return false;
        }
        return true;
    };

    let renderTableData = function(table, data, fields){
        if(typeof data === 'undefined' ||
            Object.keys(data).length == 0) {
            return;
        }

        let body = table.querySelector('tbody');
        let empty = body.querySelector('.empty');
        if(empty) {
            body.removeChild(empty);
        }
        data.forEach((item) => {
            let row = document.createElement('tr');
            fields.forEach((field) => {
                let cell = document.createElement('td');
                cell.innerHTML = item[field] ? item[field] : '&mdash;';
                row.appendChild(cell);
            });
            body.appendChild(row);
        });
    };

    let addDelimiter = function(table){
        let count = table.querySelectorAll('th').length;
        let body = table.querySelector('tbody');
        let row = document.createElement('tr');
        let cell = document.createElement('td');
        cell.classList.add('delimiter');
        cell.colSpan = count;
        row.appendChild(cell);
        body.appendChild(row);
    };

    let renderVoiceCalls = function(data) {
        let table = document.querySelector(".voice-calls");
        renderTableData(table, data, dataScheme.voice_call);
    };

    let renderArray = function(data, selector, scheme) {
        let table = document.querySelector(selector);
        data.forEach(array => {
            renderTableData(table, array, scheme);
            addDelimiter(table);
        });
    };

    let showSupportTime = function(data) {
        let colors = ['#79c779', '#eae537', '#ff4141'];
        document.querySelector(".support-info")
            .style.backgroundColor = colors[data[0] - 1];
        document.querySelector(".support-time")
            .textContent = data[1];
    };

    let snakeCaseToTitle = function(text) {
        let words = text.split('_');
        for(i in words) {
            words[i] = words[i].charAt(0).toUpperCase() +
                words[i].slice(1);
        }
        return words.join(' ');
    };

    let renderBilling = function(data) {
        if(Object.keys(data).length == 0) {
            return;
        }
        let body = document.querySelector('.billing > tbody');
        body.removeChild(body.querySelector('.empty'));
        for(key in data) {
            let row = document.createElement('tr');
            let property = document.createElement('td');
            property.classList.add('row-name');
            property.textContent = snakeCaseToTitle(key);
            let value = document.createElement('td');
            value.classList.add(data[key] ? 'true-value' : 'false-value');
            row.appendChild(property);
            row.appendChild(value);
            body.appendChild(row);
        }
    };

    let renderIncidents = function(data) {
        let table = document.querySelector(".incidents");
        renderTableData(table, data, dataScheme.incident);
    };

    let renderEmailCharts = function(data){
        let container = document.querySelector('.charts');

        console.log(data);

        Object.entries(data).forEach((country) => {
                var p = document.createElement('h3');
                p.innerText =`${country[0]}`
                container.appendChild(p);
                let div = document.createElement('div');
                country[1].forEach((item) => {
                    let labels = [];
                    let values = [];
                    item.forEach((sector) => {
                        labels.push(sector.provider + " (" + sector.country + ")");
                        values.push(sector.delivery_time);
                    });
                    let canvas = document.createElement('canvas');
                    let ctx = canvas.getContext('2d');
                    let chart = new Chart(ctx, {
                        type: 'pie',
                        data: {
                            'labels': labels,
                            'datasets': [
                                {
                                    'label': 'Dataset 1',
                                    'data': values,
                                    'backgroundColor': pieColors
                                }
                            ]
                        },
                        options: {
                            responsive: false,
                            plugins: {
                                legend: {
                                    position: 'top'
                                }
                            }
                        }
                    });
                    div.appendChild(canvas);
                });
                container.appendChild(div);
            });
    };

    let handleResponse = async function(response){
        let json = await response.json();
        if(!checkJsonScheme(json)) {
            showErrors(['JSON bad format: no status, data or error keys']);
            return;
        }

        if(json.error.length > 0) {
            showErrors([json.error]);
            return;
        }

        renderVoiceCalls(json.data.voice_call);
        renderArray(json.data.sms, ".sms", dataScheme.sms);
        renderArray(json.data.mms, ".mms", dataScheme.mms);
        showSupportTime(json.data.support);
        renderBilling(json.data.billing);
        renderIncidents(json.data.incident);

        console.log('YES');
        renderEmailCharts(json.data.email);
    };

    ready(() => {
        fetch(apiPath)
            .then(response => handleResponse(response))
            .catch(error => showErrors([error]));
    });
})();

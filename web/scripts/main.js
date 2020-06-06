/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

const baseURL = "/opencrisisline2/v1";

function handleSupportButton() {
    let number = $("#phone-number").val();
    let name = $("#name").val();
    let msg = $("#whats-going-on").val();

    sendSupportRequest(number, name, msg)
}

function sendSupportRequest(number, name, msg) {
    let data = Object();
    data.phoneNumber = number;
    data.callerName = name;
    data.message = msg;

    let dataToSend = JSON.stringify(data);
    let url = baseURL + "/support-request";
    $.ajax({
        type: 'POST',
        url: url,
        beforeSend: addHeaders,
        contentType: "application/json",
        processData: false,
        data: dataToSend,
        dataType: "text"
    }).done(function (newOb) {
        let resp = JSON.parse(newOb);
        $("#pin").text(resp.authPin)
    }).fail(function (jqXHR, textStatus, errorThrown) {
        handleError(jqXHR.statusText + " " + jqXHR.responseText);
    })
}

function addHeaders(xhr) {
    xhr.setRequestHeader('Accept', "application/json");
}

function handleError(error) {
    alert(error);
}
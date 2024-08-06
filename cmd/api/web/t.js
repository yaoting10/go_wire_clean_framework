(function (word) {
    return function encrypt(params) {
        params["uuid"] = new Date().getTime()
        let raw = JSON.stringify(params);
        console.info("request params: ", raw)
        let token = pm.environment.get("token");
        let key = token.substring(0, 16);
        var sKey = CryptoJS.enc.Utf8.parse(key);
        var sContent = CryptoJS.enc.Utf8.parse(raw);
        var encrypted = CryptoJS.AES.encrypt(sContent, sKey, {iv: sKey});
        var word = CryptoJS.enc.Base64.parse(encrypted.toString())
        var hexed = CryptoJS.enc.Hex.stringify(word)
        return hexed;
    }
})()
(function (word) {
    return function decrypt(resp) {
        let content = ""
        if (resp.code == 200) {
            content = resp.text();
        } else {
            console.error(resp.code + ": " + resp.status + ", repo: " + resp.text())
            return
        }
        let key = ""
        let token = pm.environment.get("token");
        if (token == "c4ca4238a0b923820dcc509a6f75849b-e052676e596a4662918a51e1efbe4c89" || token == "Y2FuZ2xhb3NoaXFpYW5nZ2V5YW9nZWF6a-GVsYW9kZW5nbGFvemhhb2h1Z2VodWlnZQ") {
            key = token.split("-")[1].substring(0, 16);
        } else {
            key = token.split("-")[2].substring(0, 16);
        }
        console.log("key: ", key)
        var sKey = CryptoJS.enc.Utf8.parse(key);
        var hexed = CryptoJS.enc.Hex.parse(content)
        var decrypt = CryptoJS.AES.decrypt(CryptoJS.enc.Base64.stringify(hexed), sKey, {
            iv: sKey,
        });
        let s = decrypt.toString(CryptoJS.enc.Utf8)
        console.info(resp.code + ": " + resp.status + ", repo: " + s)
        return s
    };
})()

function f() {
    const tokenUrls = ['/login/sign_in', '/login/sign_up'];
    const fixTokenUrls = ['/login/', '/systemSetting/', '/demo/'];
    var fixToken = 'cfcd208495d565ef66e7dff9f98764da-4eab9692f23d00076f62a1cbabc7fe6b';
    var CryptoJS = context.utils.CryptoJS;
    function encrypt(params) {
        params.uuid = new Date().getTime();
        let raw = JSON.stringify(params);
        console.info("request params: ", raw);
        let token = context.requestHeader.Authorization;
        token = token.split(" ")[1];
        let key = token.substring(0, 16);
        var sKey = CryptoJS.enc.Utf8.parse(key);
        var sContent = CryptoJS.enc.Utf8.parse(raw);
        var encrypted = CryptoJS.AES.encrypt(sContent, sKey, {iv: sKey});
        var word = CryptoJS.enc.Base64.parse(encrypted.toString());
        var hexed = CryptoJS.enc.Hex.stringify(word);
        return hexed;
    }

    function getToken() {
        let useFixToken = false;
        for (const u in fixTokenUrls) {
            if (context.pathname.startsWith(u)) {
                useFixToken = true;
            }
        }

        if (useFixToken) {
            return fixToken;
        } else {
            return storage.getItem("token");
        }
    }

    const token = getToken();
    if (!token) {
        console.error("接口未授权，请先登录！");
        return;
    }
    console.log("token: ", token);
    context.requestHeader.Authorization = "Bearer " + token;
    context.requestHeader['Content-Type'] = "application/x-www-form-urlencoded";

    let params = {};
    if (context.method == 'GET') {
        params = context.query;
        context.query.sign = encrypt(params);
    } else {
        params = context.requestBody;
        console.log("param:", params);
        context.requestBody.sign = encrypt(params);
    }
}
import http from "k6/http";
import { check, group } from "k6";

export default function() {
    const host = "http://localhost:8080";
    const params = {
        headers: { 
            'Content-Type': 'application/json' 
        },
    }
    const name = "jimmy";
    const email = "jimmy@example.com";
    const password = "password";

    group("register", () =>{
        let payload = JSON.stringify({
            name: name,
            email: email,
            password: password
        });
    
        let res = http.post(host + "/register", payload, params);
    
        check(res, {
            "status is 201": (r) => r.status === 201
        });
    })

    // no test because smtp server
    /* group("forgot password", () => {
        let payload = JSON.stringify({
          email: "testuser@example.com"
        });
      
        let res = http.post(host + "/forgot-password", payload, params);
      
        check(res, {
          "status is 200": (r) => r.status === 200,
          "response body is not empty": (r) => r.body.length > 0
        });
    });

    group("reset password", () =>{
        let token = "your_reset_token"; // replace with an actual token
        let payload = JSON.stringify({
            password: "new_password"
        });
    
        let res = http.post(host + "/reset-password?token=" + token, payload, params);
    
        check(res, {
            'status is 200': (r) => r.status === 200,
            'response body is not empty': (r) => r.body.length > 0,
        });
    }) */

    let token;
    let userID;
    group("login", () =>{
        let payload = JSON.stringify({
            email: email,
            password: password
        });
    
        let res = http.post(host + "/login", payload, params);
    
        check(res, {
            "status is 200": (r) => r.status === 200
        });
        

        token = res.headers["Set-Cookie"];
        let resBody = JSON.parse(res.body);
        userID = resBody.userID;
    })


    let tokenedParams = params;
    tokenedParams.headers["Cookie"] = token;

    group("view profile", () => {
        let res = http.get(`${host}/users/${userID}`, tokenedParams)
    
        check(res, {
            "status is 200": (r) => r.status === 200,
            "response body is not empty": (r) => r.body.length > 0,
        })
    })

    group("update profile", () => {
        let payload = JSON.stringify({
            "name": "zeka"
        });
    
        let res = http.patch(`${host}/${userID}`, payload, tokenedParams);
    
        check(res, {
            "status is 200": (r) => r.status === 200
        });
    });

    group('Change Password', () => {
        // Define the old and new passwords to use in the request
        const oldPassword = password;
        const newPassword = "password2";
      
        // Define the request payload
        const payload = JSON.stringify({
          old_password: oldPassword,
          new_password: newPassword,
        });
      
        // Send the request and capture the response
        const res = http.post(`${host}/users/${userID}/changepassword`, payload, tokenedParams);
      
        // Verify that the response was successful
        check(res, {
          'status is 200': (r) => r.status === 200,
          'response body is not empty': (r) => r.body.length > 0,
        });
    });


    /* const groupID = "T1"
    group("join group", () => {
        let groupID = 123;
      

        let payload = JSON.stringify({
            "group_id": groupID,
        });
      
        let res = http.post(host + "/join-group", payload, params);
      
        check(res, {
          "status is 200": (r) => r.status === 200,
          "response body is not empty": (r) => r.body.length > 0,
        });
    })

    group("leave group", () => {
        let payload = JSON.stringify({
            'group_id': '1'
        });
        let res = http.post(`${host}/users/${userID}/leavegroup`, payload, params);

        check(res, {
            'status is 200': (r) => r.status === 200,
            'response body is not empty': (r) => r.body.length > 0,
        });
    }) */

    const friendEmail = "zeus@example.com" 
    group("register friend", () =>{
        let payload = JSON.stringify({
            name: "zeus",
            email: friendEmail,
            password: "password"
        });
    
        let res = http.post(host + "/register", payload, params);
    
        check(res, {
            "status is 201": (r) => r.status === 201
        });
    })

    group("add friend", () => {
        let payload = JSON.stringify({
            "email": friendEmail,
        })

        let res = http.post(`${host}/users/${userID}/addfriend`, payload, tokenedParams);

        check(res, {
            "status is 200": (r) => r.status === 200,
            "response body has message": (r) => JSON.parse(r.body).message === "Friend added successfully",
        });
    })

    group("remove friend", () => {
        let payload = JSON.stringify({
            "email": friendEmail,
        })

        let response = http.post(`${host}/users/${userID}/removefriend`, payload, tokenedParams);
      
        check(response, {
          "status is 200": (r) => r.status === 200,
          "response body contains message": (r) => r.body.includes("Friend removed successfully"),
        });
    })

    group("logout", () =>{
        let res = http.post(host + "/logout", tokenedParams);

        check(res, {
            'status is 200': (r) => r.status === 200
        });
    })
}

import http from "k6/http";
import { check, group } from "k6";

export default function() {
    const url = "http://localhost:8080";
    const params = {
        headers: { 
            'Content-Type': 'application/json' 
        },
    }

    group("register", () =>{
        let payload = JSON.stringify({
            name: "John Doe",
            email: "john.doe@example.com",
            password: "password123"
        });
    
        let res = http.post(url + "/register", payload, params);
    
        check(res, {
            "status is 201": (r) => r.status === 201,
            "response body is not empty": (r) => r.body.length > 0
        });
    })

    group("login", () =>{
        let payload = JSON.stringify({
            email: "john.doe@example.com",
            password: "password123"
        });
    
        let res = http.post(url + "/login", payload, params);
    
        check(res, {
            "status is 200": (r) => r.status === 200,
            "response body is not empty": (r) => r.body.length > 0
        });
    })

    group("logout", () =>{
        let res = http.get(url + "/logout", params);

        check(res, {
            'status is 200': (r) => r.status === 200,
            'response body is not empty': (r) => r.body.length > 0,
        });
    })

    group("view profile", () => {
        let userId = 123 // Replace with valid user ID
        let res = http.get(`${url}/users/${userId}`, params)
    
        check(res, {
            "status is 200": (r) => r.status === 200,
            "response body is not empty": (r) => r.body.length > 0,
        })
    })

    group("update profile", () => {
        let userId = 123
        let payload = JSON.stringify({
            "name": "Jane Doe",
            "email": "jane.doe@example.com",
            "password": "newpassword123"
        });
    
        let res = http.put(`${url}/${userId}`, payload, params);
    
        check(res, {
            "status is 200": (r) => r.status === 200,
            "response body is not empty": (r) => r.body.length > 0
        });
    });

    group('Change Password', () => {
        // Define the user ID to use in the request
        const userID = 1234;
      
        // Define the old and new passwords to use in the request
        const oldPassword = 'password123';
        const newPassword = 'newpassword456';
      
        // Define the request payload
        const payload = JSON.stringify({
          old_password: oldPassword,
          new_password: newPassword,
        });
      
        // Send the request and capture the response
        const res = http.post(`${url}/change_password/${userID}`, payload, params);
      
        // Verify that the response was successful
        check(res, {
          'status is 200': (r) => r.status === 200,
          'response body is not empty': (r) => r.body.length > 0,
        });
    });

    group("forgot password", () => {
        let payload = JSON.stringify({
          email: "testuser@example.com"
        });
      
        let res = http.post(url + "/forgot-password", payload, params);
      
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
    
        let res = http.post(url + "/reset-password?token=" + token, payload, params);
    
        check(res, {
            'status is 200': (r) => r.status === 200,
            'response body is not empty': (r) => r.body.length > 0,
        });
    })

    group("join group", () => {
        let authToken = "<insert your auth token here>";
        let groupID = 123;
      

        let payload = JSON.stringify({
            "group_id": groupID,
        });
      
        let res = http.post(url + "/join-group", payload, params);
      
        check(res, {
          "status is 200": (r) => r.status === 200,
          "response body is not empty": (r) => r.body.length > 0,
        });
    })

    group("leave group", () => {
        let payload = JSON.stringify({
            'group_id': '1'
        });
        let res = http.post(`${url}/users/${userID}/leavegroup`, payload, params);

        check(res, {
            'status is 200': (r) => r.status === 200,
            'response body is not empty': (r) => r.body.length > 0,
        });
    })

    group("add friend", () => {
        const userID = 123; // replace with a valid user ID
        const friendID = 456; // replace with a valid friend ID
        const url = "http://localhost:8080/friends";

        let res = http.post(url, { friend_id: friendID }, { headers: { id: userID } });

        check(res, {
            "status is 200": (r) => r.status === 200,
            "response body has message": (r) => JSON.parse(r.body).message === "Friend added successfully",
        });
    })

    group("remove friend", () => {
        let userID = 1;
        let friendID = 2;
        
        let response = http.post(
          "http://localhost:8080/remove-friend",
          { 
            friend_id: friendID 
          },
          {
            headers: {
              "Authorization": `Bearer ${token}`,
              "Content-Type": "application/x-www-form-urlencoded"
            }
          }
        );
      
        check(response, {
          "status is 200": (r) => r.status === 200,
          "response body contains message": (r) => r.body.includes("Friend removed successfully"),
        });
    })
}

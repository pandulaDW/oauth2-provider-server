// const url = new URL(document.URL)
// const authCode = url.searchParams.get("code")
//
// const tokenData = {
//     grant_type: 'authorization_code',
//     code: authCode,
//     client_id: '000000',
//     client_secret: '999999',
//     redirect_uri: 'http://localhost:9096/callback',
// };
//
// const requestOptions = {
//     method: 'POST',
//     headers: {'Content-Type': 'application/x-www-form-urlencoded'},
//     body: new URLSearchParams(tokenData),
// };
//
// fetch("http://localhost:9096/token", requestOptions)
//     .then(response => {
//         if (!response.ok) {
//             throw new Error('Network response was not ok');
//         }
//         return response.json();
//     }).then(data => {
//     console.log('Token data:', data);
// })
//     .catch(error => {
//         console.error('There was a problem with the fetch operation:', error);
//     })

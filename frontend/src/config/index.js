export default {
    apiUrl: process.env.NODE_ENV === "production" ? window.location.host : "localhost:8787",
};

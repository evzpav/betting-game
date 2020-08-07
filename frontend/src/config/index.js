export default {
    apiUrl: process.env.NODE_ENV === "production" ? window.location.origin : "localhost:8787",
};

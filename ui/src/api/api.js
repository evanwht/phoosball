import axios from "axios";

export default axios.create({
    baseURL: 'http://localhost:3032/api/'
})
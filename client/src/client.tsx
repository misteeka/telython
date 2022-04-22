import { HttpError } from './http/error';


export default class Auth {
    static async CheckPassword(username: string, password: string): Promise<HttpError> {
        const response = await fetch(`http://127.0.0.1:8001/auth/checkPassword?u=${username}&p=${password}`, {
            method: 'GET',
            mode: 'cors',
            headers: {
                Accept: 'application/json',
            }
        });
        return new Promise<HttpError>(function (resolve, reject) {
            setTimeout(() => reject(new Error("HTTPTIMEOUT")), 1000);
            const json: any = response.json()
            let error: HttpError = json
            resolve(error)
        });
    }

}
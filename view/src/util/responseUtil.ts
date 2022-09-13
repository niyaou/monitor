export interface ICommonResponse {
    code: number;
    msg: string;
    data: any;
}

export enum CommonResponseCode {
    success = 0,
    error = 1,
}

export class CommonResponse implements ICommonResponse {
    code: number;
    msg: string;
    data: any;

    constructor(code: number, msg: string, data: any) {
        this.code = code
        this.msg = msg
        this.data = data
    }
}

export default class CommonResponseUtil {
    static success(data: any): CommonResponse {
        return new CommonResponse(CommonResponseCode.success, null, data)
    }
    static error(code: number, msg: string, data: any): CommonResponse {
        return new CommonResponse(code, msg, data)
    }
}


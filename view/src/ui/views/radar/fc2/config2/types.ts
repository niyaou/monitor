export type IForm = {
    formName: string;
    data: object;
    isSubmitStatus: boolean;
    setFailed(): void;
    setValue(formName: string, value: object): void;
    registForm(formName: string): void;
    directlySetValue(value: object): void;
};
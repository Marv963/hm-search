export interface Declaration {
  name: string;
  url: string;
}

export interface Option {
  name: string;
  title: string;
  description: string;
  type: string;
  default: string;
  example: string;
  readOnly: boolean;
  loc: string[];
  declarations: Declaration[];
}

export interface OptionsData {
  last_update: string;
  count: number;
  options: Option[];
}

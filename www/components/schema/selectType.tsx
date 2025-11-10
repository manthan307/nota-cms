import * as React from "react";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { InputTypes } from "@/lib/types";

export type SelectTypeProps = {
  name: string;
  defaultValue?: string;
  disable?: boolean;
};
export function SelectType({ name, defaultValue, disable }: SelectTypeProps) {
  const [value, setValue] = React.useState(defaultValue || "");

  return (
    <Select value={value} onValueChange={setValue} disabled={disable}>
      <SelectTrigger className="w-[100px]">
        <SelectValue placeholder="Type" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Type</SelectLabel>
          {InputTypes.map((value, index) => (
            <SelectItem value={value} key={index}>
              {value}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>

      {/* Hidden input for form submission */}
      <input type="hidden" name={name} value={value} />
    </Select>
  );
}

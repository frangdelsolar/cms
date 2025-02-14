import Form from "@rjsf/mui";
import validator from "@rjsf/validator-ajv8";
import { useEffect, useState } from "react";
import { useAppDispatch, useAppSelector } from "../../../../store/Hooks";
import {
  selectFormData,
  setFormData,
  selectFormErrors,
  selectFormInitialized,
  selectFormSaving,
  setFormInitialized,
  setFormSaving,
  clearForm,
} from "../../../../store/FormSlice";

const ModelForm = ({ schema, data, submitHandler }) => {
  const [formattedErrors, setFormattedErrors] = useState([]);

  // Clean up schema and data
  const cleanedSchema = cleanSchema(schema);
  const cleanedData = cleanData(data, cleanedSchema.removedKeys);

  // Redux hooks
  const dispatch = useAppDispatch();
  const formData = useAppSelector(selectFormData);
  const formErrors = useAppSelector(selectFormErrors);
  const formSaving = useAppSelector(selectFormSaving);
  const formInitialized = useAppSelector(selectFormInitialized);

  // Initialize the form state with cleaned data
  useEffect(() => {
    dispatch(setFormData(cleanedData));
    dispatch(setFormInitialized(true));

    return () => {
      dispatch(clearForm()); // Clear form data on unmount
    };
  }, [data, schema]);

  useEffect(() => {
    setFormattedErrors(FormatErrors());
  }, [formErrors]);

  useEffect(() => {
    if (formSaving && typeof submitHandler === "function") {
      submitHandler(formData);
      dispatch(setFormSaving(false));
    }
  }, [formSaving]);

  const handleOnChange = (e) => {
    dispatch(setFormData(e.formData));
  };

  const FormatErrors = () => {
    let output = {};
    for (const error of formErrors) {
      let fieldKey = error["Field"];
      let fieldValue = error["Error"];

      if (Object.keys(output).includes(fieldKey)) {
        output[fieldKey].__errors.push(fieldValue);
      } else {
        output[fieldKey] = { __errors: [fieldValue] };
      }
    }
    return output;
  };

  return (
    <>
      {formInitialized && (
        <Form
          uiSchema={GetUISchema(data)}
          schema={cleanedSchema.def} // The cleaned schema to render the form
          formData={formData} // Bind the form data from Redux state
          onChange={handleOnChange} // Handle form data changes
          validator={validator} // Validator for form
          templates={GetTemplates()} // Custom submit button template
          // widgets={ GetWidgets() }
          extraErrors={formattedErrors}
        />
      )}
    </>
  );
};

function SubmitButton() {
  return <></>; // Placeholder submit button (can be customized or removed)
}

function GetTemplates() {
  return {
    ButtonTemplates: { SubmitButton },
  };
}

function GetUISchema(data) {
  let output = {};

  if (!data) {
    return output;
  }

  // Set textarea widget for long text fields
  Object.keys(data).forEach((key) => {
    let value = data[key];
    if (value && value.length > 60) {
      output[key] = { "ui:widget": "textarea" };
    }
  });

  return output;
}

/**
 * Copies a schema at a given path. This is used to safely extract a part of the schema.
 *
 * @param {Object} schema - The schema from which to extract a portion.
 * @param {Array} path - The path to the part of the schema to copy.
 * @returns {Object|undefined} - The copied portion of the schema or undefined if not found.
 */
function copySchemaAtPath(schema, path) {
  if (!schema || !path || !Array.isArray(path) || path.length === 0) {
    return undefined; // Return undefined if path is invalid
  }

  let current = schema;
  for (const key of path) {
    if (current && typeof current === "object" && key in current) {
      current = current[key]; // Drill down the path
    } else {
      return undefined; // Return undefined if path doesn't exist
    }
  }

  if (current && typeof current === "object") {
    return JSON.parse(JSON.stringify(current)); // Deep copy to avoid mutating the original schema
  } else {
    return undefined; // Return undefined if the final object is not valid
  }
}

/**
 * Cleans the schema by removing foreign references and unused keys.
 *
 * @param {Object} schema - The schema to clean.
 * @returns {Object} - The cleaned schema and removed keys.
 */
function cleanSchema(schema) {
  // Extract path from the $ref to locate the schema definition
  let path = schema["$ref"].split("/").slice(1);
  let def = copySchemaAtPath(schema, path); // Get the definition from the schema

  let removedKeys = []; // Keep track of removed keys for data cleaning

  // Remove foreign object references from the properties
  for (const key in def.properties) {
    if ("$ref" in def.properties[key]) {
      delete def.properties[key]; // Remove the property with foreign reference

      // Remove key from the required properties list
      def.required = def.required.filter((item) => item !== key);

      removedKeys.push(key); // Track the removed key
    }
  }

  // TODO: Maybe we can handle this from the backend. But for now...
  delete def.required;

  def.description =
    "SystemData fields (ID, CreatedAt, UpdatedAt, CreatedById and UpdatedById) will be automatically populated by the server.";

  return { def, removedKeys }; // Return the cleaned schema and removed keys
}

/**
 * Cleans the form data by removing keys that have been removed from the schema
 * and formatting dates.
 *
 * @param {Object} data - The form data to clean.
 * @param {Array} removedKeys - The list of keys to be removed.
 * @returns {Object} - The cleaned data.
 */
function cleanData(data, removedKeys) {
  let output = {};

  // Remove the specified keys from data
  for (const key in data) {
    let remove = removedKeys.find((item) => item === key) || key === "id"; // Skip removed or "id" keys
    if (!remove) {
      output[key] = data[key]; // Include the key if it's not removed

      // If the value is a date, convert it to a string
      if (data[key] instanceof Date) {
        output[key] = data[key].toISOString();
      }
    }
  }
  return output; // Return the cleaned data
}

export default ModelForm;

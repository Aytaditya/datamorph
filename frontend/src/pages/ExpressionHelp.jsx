import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Code, Copy, CheckCircle } from 'lucide-react';
import toast from 'react-hot-toast';

const ExpressionHelp = () => {
  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    toast.success('Copied to clipboard!');
  };

  const examples = [
    {
      title: "Basic Value Access",
      description: "Access the source value directly",
      expression: "value",
      example: "If source is 'John', output will be 'John'"
    },
    {
      title: "String Manipulation",
      description: "Convert to uppercase",
      expression: "toUpper(value)",
      example: "If source is 'john', output will be 'JOHN'"
    },
    {
      title: "String Concatenation",
      description: "Combine multiple values",
      expression: "value + ' - Customer'",
      example: "If source is 'John', output will be 'John - Customer'"
    },
    {
      title: "Conditional Logic",
      description: "Use if-else conditions",
      expression: "value != '' ? toUpper(value) : 'UNKNOWN'",
      example: "If source is empty, output 'UNKNOWN', otherwise uppercase the value"
    },
    {
      title: "Number Operations",
      description: "Mathematical calculations",
      expression: "value * 1.2",
      example: "If source is 100, output will be 120"
    },
    {
      title: "Date Functions",
      description: "Get current date",
      expression: "getCurrentDate()",
      example: "Output will be current date in YYYY-MM-DD format"
    },
    {
      title: "Default Values",
      description: "Provide fallback for missing data",
      expression: "value || 'Default Value'",
      example: "If source is empty/null, output 'Default Value'"
    },
    {
      title: "Complex Logic",
      description: "Multiple conditions and operations",
      expression: "len(value) > 5 ? toUpper(value) : toLower(value) + '_short'",
      example: "If value length > 5, uppercase it; otherwise lowercase and add '_short'"
    }
  ];

  const functions = [
    {
      name: "toUpper(str)",
      description: "Convert string to uppercase",
      example: "toUpper('hello') → 'HELLO'"
    },
    {
      name: "toLower(str)",
      description: "Convert string to lowercase", 
      example: "toLower('HELLO') → 'hello'"
    },
    {
      name: "len(str)",
      description: "Get length of string",
      example: "len('hello') → 5"
    },
    {
      name: "substr(str, start, length)",
      description: "Extract substring",
      example: "substr('hello', 1, 3) → 'ell'"
    },
    {
      name: "contains(str, substring)",
      description: "Check if string contains substring",
      example: "contains('hello world', 'world') → true"
    },
    {
      name: "replace(str, old, new)",
      description: "Replace text in string",
      example: "replace('hello world', 'world', 'universe') → 'hello universe'"
    },
    {
      name: "trim(str)",
      description: "Remove whitespace from start and end",
      example: "trim('  hello  ') → 'hello'"
    },
    {
      name: "getCurrentDate()",
      description: "Get current date in YYYY-MM-DD format",
      example: "getCurrentDate() → '2025-07-10'"
    },
    {
      name: "formatDate(date, format)",
      description: "Format date string",
      example: "formatDate('2025-07-10', 'DD/MM/YYYY') → '10/07/2025'"
    }
  ];

  const operators = [
    { symbol: "+", description: "Addition or string concatenation", example: "5 + 3 → 8, 'Hello' + ' World' → 'Hello World'" },
    { symbol: "-", description: "Subtraction", example: "10 - 3 → 7" },
    { symbol: "*", description: "Multiplication", example: "4 * 5 → 20" },
    { symbol: "/", description: "Division", example: "15 / 3 → 5" },
    { symbol: "==", description: "Equality comparison", example: "value == 'test' → true/false" },
    { symbol: "!=", description: "Not equal comparison", example: "value != '' → true/false" },
    { symbol: ">", description: "Greater than", example: "len(value) > 5 → true/false" },
    { symbol: "<", description: "Less than", example: "value < 100 → true/false" },
    { symbol: ">=", description: "Greater than or equal", example: "value >= 18 → true/false" },
    { symbol: "<=", description: "Less than or equal", example: "value <= 65 → true/false" },
    { symbol: "&&", description: "Logical AND", example: "value > 0 && value < 100 → true/false" },
    { symbol: "||", description: "Logical OR", example: "value == 'admin' || value == 'user' → true/false" },
    { symbol: "?:", description: "Ternary operator (if-else)", example: "value > 18 ? 'adult' : 'minor'" }
  ];

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">
            Expression Help & Reference
          </h1>
          <p className="text-lg text-gray-600 max-w-3xl mx-auto">
            Complete guide for writing transformation expressions in mapping rules. 
            Use these expressions to transform data during the mapping process.
          </p>
        </div>

        {/* Quick Start */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center">
              <Code className="h-5 w-5 mr-2" />
              Quick Start
            </CardTitle>
            <CardDescription>
              Basic concepts for writing expressions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <h4 className="font-semibold mb-2">Key Points:</h4>
                <ul className="space-y-1 text-sm text-gray-600">
                  <li>• Use <code className="bg-gray-100 px-1 rounded">value</code> to access the source field value</li>
                  <li>• Expressions are JavaScript-like but simplified</li>
                  <li>• String literals must be in single quotes: <code className="bg-gray-100 px-1 rounded">'text'</code></li>
                  <li>• Use operators for logic: <code className="bg-gray-100 px-1 rounded">&&</code>, <code className="bg-gray-100 px-1 rounded">||</code>, <code className="bg-gray-100 px-1 rounded">?:</code></li>
                  <li>• Functions are available for common operations</li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold mb-2">Simple Example:</h4>
                <div className="bg-gray-100 p-3 rounded text-sm">
                  <div><strong>Source:</strong> <code>"john doe"</code></div>
                  <div><strong>Expression:</strong> <code>toUpper(value)</code></div>
                  <div><strong>Result:</strong> <code>"JOHN DOE"</code></div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Examples */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Common Expression Examples</CardTitle>
            <CardDescription>
              Ready-to-use expressions for typical transformation scenarios
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid lg:grid-cols-2 gap-4">
              {examples.map((example, index) => (
                <div key={index} className="border rounded-lg p-4">
                  <div className="flex justify-between items-start mb-2">
                    <h4 className="font-semibold text-sm">{example.title}</h4>
                    <button
                      onClick={() => copyToClipboard(example.expression)}
                      className="text-gray-400 hover:text-gray-600"
                      title="Copy expression"
                    >
                      <Copy className="h-4 w-4" />
                    </button>
                  </div>
                  <p className="text-xs text-gray-600 mb-2">{example.description}</p>
                  <div className="bg-gray-50 p-2 rounded text-xs font-mono mb-2">
                    {example.expression}
                  </div>
                  <p className="text-xs text-green-600">{example.example}</p>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Functions Reference */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Available Functions</CardTitle>
            <CardDescription>
              Built-in functions you can use in expressions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid lg:grid-cols-2 gap-4">
              {functions.map((func, index) => (
                <div key={index} className="border rounded-lg p-4">
                  <div className="flex justify-between items-start mb-2">
                    <h4 className="font-semibold text-sm font-mono">{func.name}</h4>
                    <button
                      onClick={() => copyToClipboard(func.name)}
                      className="text-gray-400 hover:text-gray-600"
                      title="Copy function"
                    >
                      <Copy className="h-4 w-4" />
                    </button>
                  </div>
                  <p className="text-xs text-gray-600 mb-2">{func.description}</p>
                  <div className="bg-green-50 p-2 rounded text-xs font-mono text-green-700">
                    {func.example}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Operators Reference */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Operators Reference</CardTitle>
            <CardDescription>
              Available operators for building expressions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {operators.map((op, index) => (
                <div key={index} className="flex items-start space-x-4 border-b pb-3 last:border-b-0">
                  <Badge variant="outline" className="font-mono">{op.symbol}</Badge>
                  <div className="flex-1">
                    <p className="text-sm font-medium">{op.description}</p>
                    <p className="text-xs text-gray-600 mt-1">{op.example}</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Best Practices */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <CheckCircle className="h-5 w-5 mr-2 text-green-600" />
              Best Practices
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <h4 className="font-semibold text-green-600 mb-2">✅ Do:</h4>
                <ul className="space-y-1 text-sm">
                  <li>• Use meaningful variable names in complex expressions</li>
                  <li>• Test expressions with sample data first</li>
                  <li>• Handle empty/null values with fallbacks</li>
                  <li>• Keep expressions readable and simple when possible</li>
                  <li>• Use parentheses to clarify operator precedence</li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold text-red-600 mb-2">❌ Don't:</h4>
                <ul className="space-y-1 text-sm">
                  <li>• Create overly complex expressions in one line</li>
                  <li>• Forget to handle missing data cases</li>
                  <li>• Use double quotes for strings (use single quotes)</li>
                  <li>• Assume data types without validation</li>
                  <li>• Nest too many conditional operations</li>
                </ul>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Footer */}
        <div className="text-center mt-8 text-gray-500 text-sm">
          <p>Need help? Contact support or check the transformation logs for debugging.</p>
        </div>
      </div>
    </div>
  );
};

export default ExpressionHelp;

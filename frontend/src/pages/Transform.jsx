import React, { useState, useEffect } from 'react';
import { Play, Download, Upload, AlertCircle, CheckCircle } from 'lucide-react';
import { clientsAPI, transformAPI } from '../services/api';
import toast from 'react-hot-toast';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';

const Transform = () => {
  const [clients, setClients] = useState([]);
  const [selectedClient, setSelectedClient] = useState('');
  const [inputData, setInputData] = useState('{\n  "example": {\n    "data": "value"\n  }\n}');
  const [outputData, setOutputData] = useState('');
  const [loading, setLoading] = useState(false);
  const [loadingClients, setLoadingClients] = useState(true);

  useEffect(() => {
    loadClients();
  }, []);

  const loadClients = async () => {
    try {
      const data = await clientsAPI.getAll();
      setClients(data);
    } catch (error) {
      toast.error('Failed to load clients');
    } finally {
      setLoadingClients(false);
    }
  };

  const handleTransform = async () => {
    if (!selectedClient) {
      toast.error('Please select a client');
      return;
    }

    if (!inputData.trim()) {
      toast.error('Please provide input data');
      return;
    }

    setLoading(true);
    try {
      const parsedInput = JSON.parse(inputData);
      const result = await transformAPI.transform(selectedClient, parsedInput);
      
      // Simply display whatever the backend returns
      setOutputData(JSON.stringify(result, null, 2));
      
      if (result.warnings && result.warnings.missingRequiredFields) {
        toast.warning(
          `Transformation completed with warnings: ${result.warnings.missingRequiredFields.length} missing required fields`
        );
      } else {
        toast.success('Transformation completed successfully');
      }
    } catch (error) {
      console.error('Transformation error:', error);
      if (error.name === 'SyntaxError') {
        toast.error('Invalid JSON format in input data');
      } else if (error.response?.data?.error) {
        toast.error(`Transformation failed: ${error.response.data.error}`);
      } else {
        toast.error('Transformation failed');
      }
      setOutputData('');
    } finally {
      setLoading(false);
    }
  };

  const handleLoadSampleData = () => {
    const sampleData = {
      "applicantDetails": [
        {
          "loanId": 5571,
          "accountNo": "000005571",
          "entityType": "Applicant",
          "entityName": "John Doe",
          "mobileNo": "9876543210",
          "fatherName": "Robert Doe",
          "qualification": "GRADUATE",
          "gender": "MALE",
          "relationship": "Self",
          "profession": "SALARIED",
          "dob": "15-January-1990"
        }
      ],
      "applicantKyc": [
        {
          "loanId": 5571,
          "entityType": "Applicant",
          "entityName": "John Doe",
          "panNumber": "ABCDE1234F",
          "aadhaarNumber": "XXXXXXXXXXXX1234"
        }
      ]
    };
    
    setInputData(JSON.stringify(sampleData, null, 2));
    toast.success('Sample data loaded');
  };

  const handleDownloadOutput = () => {
    if (!outputData) {
      toast.error('No output data to download');
      return;
    }

    const blob = new Blob([outputData], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `transformed-data-${Date.now()}.json`;
    a.click();
    URL.revokeObjectURL(url);
    toast.success('Output data downloaded');
  };

  if (loadingClients) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="border-b border-gray-200 pb-5">
        <h1 className="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
          JSON Transformation
        </h1>
        <p className="mt-1 text-sm text-gray-500">
          Test data transformations using client-specific mapping rules
        </p>
      </div>

      {/* Client Selection */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Client Selection</CardTitle>
              <CardDescription>Choose a client to test data transformations</CardDescription>
            </div>
            <Button
              onClick={handleLoadSampleData}
              variant="outline"
              className="flex items-center"
            >
              <Upload className="h-4 w-4 mr-2" />
              Load Sample Data
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {clients.length === 0 ? (
            <div className="text-center py-8">
              <AlertCircle className="mx-auto h-12 w-12 text-gray-400" />
              <h3 className="mt-2 text-sm font-medium text-gray-900">No clients available</h3>
              <p className="mt-1 text-sm text-gray-500">
                Create a client first to test transformations.
              </p>
            </div>
          ) : (
            <div className="max-w-sm">
              <Select
                value={selectedClient}
                onValueChange={(value) => setSelectedClient(value)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a client..." />
                </SelectTrigger>
                <SelectContent>
                  {clients.map((client) => (
                    <SelectItem key={client.id} value={client.id.toString()}>
                      {client.name} (ID: {client.id})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Transformation Interface */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Input Panel */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Input JSON</CardTitle>
              <CardDescription>Source Data</CardDescription>
            </div>
          </CardHeader>
          <CardContent>
            <Textarea
              value={inputData}
              onChange={(e) => setInputData(e.target.value)}
              className="font-mono text-sm min-h-[500px]"
              placeholder="Enter your JSON data here..."
            />
            <div className="mt-4 flex justify-center">
              <Button
                onClick={handleTransform}
                disabled={loading || !selectedClient}
                className="flex items-center"
              >
                {loading ? (
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                ) : (
                  <Play className="h-4 w-4 mr-2" />
                )}
                {loading ? 'Transforming...' : 'Transform Data'}
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Output Panel */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Output JSON</CardTitle>
              <div className="flex items-center space-x-2">
                {outputData && (
                  <>
                    <CheckCircle className="h-5 w-5 text-green-500" />
                    <CardDescription>Transformed Data</CardDescription>
                    <Button
                      onClick={handleDownloadOutput}
                      variant="outline"
                      size="sm"
                      className="flex items-center"
                    >
                      <Download className="h-4 w-4 mr-1" />
                      Download
                    </Button>
                  </>
                )}
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <Textarea
              value={outputData}
              readOnly
              className="font-mono text-sm bg-gray-50 min-h-[500px]"
              placeholder="Transformed data will appear here..."
            />
            {outputData && (
              <div className="mt-4 p-3 bg-green-50 border border-green-200 rounded-lg">
                <div className="flex items-center">
                  <CheckCircle className="h-5 w-5 text-green-400 mr-2" />
                  <span className="text-sm font-medium text-green-800">
                    Transformation completed successfully
                  </span>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Instructions */}
      <Card className="bg-blue-50 border-blue-200">
        <CardContent className="pt-6">
          <div className="flex items-start">
            <AlertCircle className="h-5 w-5 text-blue-400 mt-0.5 mr-3" />
            <div>
              <h3 className="text-sm font-medium text-blue-800">How to use</h3>
              <div className="mt-2 text-sm text-blue-700">
                <ol className="list-decimal list-inside space-y-1">
                  <li>Select a client that has mapping rules configured</li>
                  <li>Enter or paste your JSON data in the input panel</li>
                  <li>Click "Transform Data" to apply the mapping rules</li>
                  <li>Review the transformed output and download if needed</li>
                </ol>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Transform;

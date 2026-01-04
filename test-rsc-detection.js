// Test RSC detection logic
const testCases = [
    { url: '/dashboard?_rsc=abc', expected: true, description: 'URL with _rsc param' },
    { url: '/dashboard', expected: false, description: 'URL without _rsc param' },
    { url: '/api/test', expected: false, description: 'API route' },
    { url: '/dashboard?_rsc=', expected: true, description: 'URL with empty _rsc param' },
];

console.log('Testing RSC detection logic:');
testCases.forEach(test => {
    const hasRscParam = test.url.includes('_rsc=');
    const passed = hasRscParam === test.expected;
    console.log(`${passed ? '✅' : '❌'} ${test.description}: ${test.url} (has _rsc=${hasRscParam}, expected=${test.expected})`);
});

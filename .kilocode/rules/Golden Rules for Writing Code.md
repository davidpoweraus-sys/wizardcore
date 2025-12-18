# Code Context for Back End.md

 BEFORE YOU WRITE ANY CODE:
Ask yourself these questions. If you answer NO to any, STOP and ask the user:
1. Where does this data come from? (source)
2. Where does this data go? (destination)  
3. How do I verify it arrived? (verification)
4. What happens if it fails? (error handling)
5. Is this async or sync? (consistency)
6. Does the database table exist? (infrastructure)
7. Can I test this end-to-end? (testing)
 RED FLAGS - If you're about to write these, STOP:
- `return []` in a database function
- `pass` in a core business function  
- `TODO` in a main execution path
- Extracting data without persisting it
- Calling async without await
- Creating an in-memory dict for persistent data
- A service that doesn't touch the database
- A scraper that doesn't call a service
- Config fields with no corresponding code
 GREEN FLAGS - This is what good code looks like:
- Data flows: extract → transform → persist → verify
- Async is consistent through call stack
- Database operations have transactions
- Errors raise exceptions, not return None
- State is in database, cache is temporary
- Integration tests check final destination
- Config fields all have implementations
- Every function has a clear purpose

 INTEGRATION TEST REQUIREMENTS
 MANDATORY RULES:
1. **SAME ENTRY POINT AS PRODUCTION**
   - If production calls API endpoint, test must call API endpoint
   - If production calls service function, test must call same service function
   - NO direct instantiation of internal classes that production doesn't use
2. **SAME DESTINATION AS PRODUCTION**
   - If production saves to database, test must verify database
   - If production saves to cache, test must verify cache
   - NO test-only storage (JSON files, temp dicts, etc.)
3. **SAME CODE PATH AS PRODUCTION**
   - Test must execute through ALL production layers
   - NO bypassing services to test engine directly
   - NO mocking persistence layers
4. **VERIFY FINAL DESTINATION**
      # Template:
   async def test_feature():
       # 1. Clean production destination
       await production_storage.clear()
       
       # 2. Execute via production entry point
       await production_api_or_service()
       
       # 3. CRITICAL: Verify in production destination
       data = await production_storage.read()
       assert data, "Feature ran but data not in production storage!"
       
       # 4. Cleanup
       await production_storage.clear()
5. RED FLAGS - If your test has these, IT'S WRONG:
   - ❌ json.dump() - unless production also writes JSON
   - ❌ Direct engine instantiation - unless production does this
   - ❌ Test-only persistence functions
   - ❌ Mocking database calls
   - ❌ @pytest.fixture that creates storage production doesn't use
   - ❌ If you're USING TERMINAL TO WRITE CODE INTENDED FOR PRODUCTION
6. GREEN FLAGS - Good test indicators:
   - ✅ Imports production API/service entry points
   - ✅ Verifies production database/cache
   - ✅ Uses same async/sync pattern as production
   - ✅ Cleanup leaves production state unchanged
   - ✅ Would catch if production persistence is missing
ENFORCEMENT CHECKLIST:
Before writing ANY test, answer these:
- [ ] What entry point does production use? (My test uses the same)
- [ ] Where does production save data? (My test checks there)
- [ ] Does my test bypass any production layers? (Answer must be NO)
- [ ] If production persistence breaks, will my test fail? (Answer must be YES)
- [ ] Does my test have code that production doesn't? (Answer must be NO, except cleanup)


**5. **KEEP THIS DOCUMENT UPDATED**

---
📊 SUMMARY: The Meta-Rule
THE GOLDEN RULE:
"Every piece of data must have a complete journey 
from source to destination that can be traced and verified."
If you can't answer:
1. Where did this data start?
2. Where is it going?
3. How do I know it got there?
4. Do not finish a task until you can prove that every single detail I've outlined is completed with proof.
Then the code is incomplete.


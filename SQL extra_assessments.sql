-- Populate the Assessments table
INSERT INTO Assessments (name, description)
VALUES 
    ('Cardiovascular Health Check', 'Assessment to evaluate heart disease risk factors.'),
    ('Mental Wellness Assessment', 'Assessment to evaluate stress and anxiety levels.');

-- Populate the Questions table
INSERT INTO Questions (assessment_id, question_text, type)
VALUES
   
    -- Questions for Cardiovascular Health Check
    (2, 'How often do you consume fried food?', 'mcq'),
    (2, 'Do you have a family history of heart disease?', 'mcq'),
    (2, 'Do you engage in regular physical activity?', 'mcq'),
    (2, 'How often do you check your blood pressure?', 'mcq'),

    -- Questions for Mental Wellness Assessment
    (3, 'How often do you feel stressed?', 'mcq'),
    (3, 'Do you have trouble sleeping at night?', 'mcq'),
    (3, 'Do you engage in relaxation activities?', 'mcq'),
    (3, 'How often do you feel overwhelmed?', 'mcq');

-- Populate the Options table
INSERT INTO Options (assessment_id, question_id, option_text, risk_value)
VALUES
 
    -- Cardiovascular Health Check Options
    (2, 5, 'Rarely', 0), (2, 5, '1-2 times a week', 1), (2, 5, '3-4 times a week', 2), (2, 5, 'Daily', 3),
    (2, 6, 'Yes', 3), (2, 6, 'No', 0),
    (2, 7, 'Yes', 0), (2, 7, 'No', 3),
    (2, 8, 'Monthly', 0), (2, 8, 'Few times a year', 1), (2, 8, 'Rarely', 2), (2, 8, 'Never', 3),

    -- Mental Wellness Assessment Options
    (3, 9, 'Rarely', 0), (3, 9, 'Sometimes', 1), (3, 9, 'Often', 2), (3, 9, 'Always', 3),
    (3, 10, 'Yes', 3), (3, 10, 'No', 0),
    (3, 11, 'Yes', 0), (3, 11, 'No', 3),
    (3, 12, 'Rarely', 0), (3, 12, 'Sometimes', 1), (3, 12, 'Often', 2), (3, 12, 'Always', 3);

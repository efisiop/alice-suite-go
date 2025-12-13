# ALICE READER APP - STRUCTURED OUTPUT SCHEMA SPECIFICATIONS

## API RESPONSE SCHEMAS

### STANDARD BASE RESPONSE WRAPPER
All API responses follow consistent schema with metadata, results, and error handling.

```json
{
  "success": true,
  "timestamp": "2025-11-21T14:30:00Z",
  "request_id": "req_abc123def456",
  "version": "v1.0",
  "data": {},
  "metadata": {
    "processing_time_ms": 45,
    "cache_hit": false,
    "rate_limit_remaining": 1000
  },
  "error": null
}
```

**Field Specifications:**
- `success`: Boolean indicating request success/failure
- `timestamp`: ISO 8601 timestamp of response generation
- `request_id`: Unique identifier for request tracking
- `version`: API semantic version
- `data`: Response payload (schema defined per endpoint)
- `metadata`: Request processing details
- `error`: Error details if `success` is false

### BOOK STRUCTURE SCHEMAS

#### Book Metadata Response
```json
{
  "success": true,
  "data": {
    "id": "book_alice_adventures_wonderland",
    "title": "Alice's Adventures in Wonderland",
    "author": "Lewis Carroll",
    "edition": "1865 First Edition",
    "total_pages": 200,
    "language": "en",
    "isbn": "978-1234567890",
    "cover_image_url": "/assets/covers/alice.jpg",
    "description": "Alice falls down a rabbit hole into a fantasy world...",
    "reading_level": "intermediate",
    "genres": ["fantasy", "childrens", "literature"],
    "publisher": "Macmillan",
    "publication_year": 1865,
    "last_updated": "2025-11-20T16:45:00Z"
  }
}
```

#### Page Response Schema
```json
{
  "success": true,
  "data": {
    "book_id": "book_alice_adventures_wonderland",
    "page_number": 24,
    "word_count": 187,
    "sections": [
      {
        "id": "section_p24_s1",
        "section_number": 1,
        "content": "There was nothing so very remarkable in that...",
        "word_count": 41,
        "glossary_terms": [
          {
            "term": "remarkable",
            "offset": 31,
            "length": 10,
            "glossary_definition_id": "def_remarkable"
          }
        ]
      },
      {
        "id": "section_p24_s2",
        "section_number": 2,
        "content": "Nor did Alice think it so very much out of the way...",
        "word_count": 44,
        "glossary_terms": []
      }
    ],
    "navigation": {
      "previous_page_available": true,
      "previous_page": 23,
      "next_page_available": true,
      "next_page": 25
    }
  }
}
```

### DICTIONARY LOOKUP SCHEMAS

#### Word Definition Response
```json
{
  "success": true,
  "data": {
    "query": "grinned",
    "found_in": "alice_glossary", // or "external_dictionary"
    "definitions": [
      {
        "source": "alice_glossary",
        "definition": "smiled broadly, especially in an unrestrained manner",
        "part_of_speech": "verb",
        "context_sentence": "The Cheshire Cat grinned from ear to ear.",
        "alice_specific": true,
        "page_reference": "Page 12, Section 3",
        "etymology": "Old English grennian 'to show the teeth'"
      }
    ],
    "alternative_sources": [
      {
        "source": "merriam_webster",
        "definition": "to draw back the lips and show the teeth",
        "part_of_speech": "verb",
        "alice_specific": false
      }
    ],
    "related_terms": ["smiled", "beamed", "leered"],
    "difficulty_level": "intermediate",
    "lookup_count": 127,
    "previous_user_lookup": "2025-11-18T09:15:00Z"
  }
}
```

#### Multi-Word Selection Response (AI Trigger)
```json
{
  "success": true,
  "data": {
    "selection": {
      "text": "curious dream-like quality",
      "word_count": 4,
      "text_type": "phrase" // 1-5 words = dictionary, 6+ = ai_assistant
    },
    "is_eligible_for_ai": false,
    "handling": "dictionary_lookup",
    "ai_trigger_threshold": 5, // words
    "current_word_count": 4,
    "needs_additional_words": 1,
    "dictionary_definitions": [
      {
        "word": "curious",
        "definitions": [/* ... */]
      },
      {
        "word": "dream-like",
        "definitions": [/* ... */]
      },
      {
        "word": "quality",
        "definitions": [/* ... */]
      }
    ]
  }
}
```

### AI ASSISTANT SCHEMAS

#### AI Explanation Request/Response
```json
// Request Schema
{
  "selection": {
    "text": "curious dream-like quality of the entire adventure",
    "word_count": 8,
    "context": {
      "current_page": 24,
      "current_section": 2,
      "book": "Alice's Adventures in Wonderland",
      "current_theme": "surreal transformation"
    }
  },
  "interaction_type": "explain", // CHAT, EXPLAIN, QUIZ, SIMPLIFY, DEFINITION
  "previous_interactions": [
    {
      "id": "ai_int_001",
      "timestamp": "2025-11-21T14:25:00Z",
      "type": "explain",
      "selection": "down the rabbit hole"
    }
  ],
  "user_progress": {
    "current_page": 24,
    "pages_read_today": 5,
    "reading_velocity_wpm": 120,
    "vocabulary_level": "intermediate"
  }
}

// Response Schema
{
  "success": true,
  "data": {
    "ai_response": {
      "id": "ai_resp_001",
      "timestamp": "2025-11-21T14:30:15Z",
      "explanation": "Lewis Carroll uses the phrase 'curious dream-like quality' to describe the surreal, illogical nature of Wonderland. Unlike regular dreams, Alice's adventures follow consistent internal logic, but defy normal reality",
      "original_selection": "curious dream-like quality of the entire adventure",
      "simplification": "The whole story feels like a strange but interesting dream",
      "key_concepts": ["surrealism", "dream logic", "narrative tone"],
      "contextual_examples": [
        "This dream-like quality is established from the moment Alice falls down the rabbit hole",
        "Wonderland operates on its own internal rules, much like dreams do"
      ],
      "reading_level_appropriate": true,
      "confidence_score": 0.87,
      "interaction_type": "explain",
      "response_time_ms": 923
    },
    "related_questions": [
      {
        "text": "How does this dream-like quality affect the story's structure?",
        "question_type": "analysis"
      },
      {
        "text": "Are there other books with similar dream-like qualities?",
        "question_type": "comparison"
      }
    ],
    "quiz_questions": [
      {
        "question": "What does 'curious dream-like quality' suggest about Wonderland's reality?",
        "options": ["Strictly logical", "Surreal and illogical", "Completely random", "Based on science"],
        "correct_answer": 1,
        "explanation": "The dream-like quality implies Wonderland follows its own unique logic, different from normal reality"
      }
    ]
  }
}
```

### USER MANAGEMENT SCHEMAS

#### Authentication Response
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user_abc123def456",
      "email": "alice.student@school.edu",
      "first_name": "Alice",
      "last_name": "Johnson",
      "role": "reader", // reader, consultant, admin
      "created_at": "2025-09-15T10:30:00Z",
      "last_login": "2025-11-21T08:15:00Z",
      "preferences": {
        "reading_font_size": "1.1rem",
        "line_height": 1.8,
        "theme": "light",
        "auto_save": true,
        "analytics_tracking": true
      }
    },
    "session": {
      "jwt_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_at": "2025-11-21T18:30:00Z",
      "refresh_token": "refresh_xyz789abc123",
      "permissions": ["read_books", "track_progress", "request_help"]
    }
  }
}
```

#### Reading Progress Response
```json
{
  "success": true,
  "data": {
    "current_position": {
      "book_id": "book_alice_adventures_wonderland",
      "page_number": 67,
      "section_number": 3,
      "position_percentage": 33.5,
      "last_accessed": "2025-11-21T14:30:00Z"
    },
    "statistics": {
      "pages_read": 67,
      "pages_remaining": 133,
      "sections_completed": 289,
      "total_pages": 200,
      "total_sections": 850,
      "completion_percentage": 33.5,
      "reading_velocity_wpm": 124,
      "estimated_time_remaining_minutes": 428
    },
    "today_progress": {
      "pages_read": 8,
      "reading_time_minutes": 35,
      "vocabulary_lookups": 12,
      "ai_interactions": 3,
      "help_requests": 0
    },
    "weekly_goals": {
      "target_pages": 50,
      "pages_read": 37,
      "target_met_percentage": 74,
      "streak_days": 5
    },
    "vocabulary_building": {
      "words_looked_up": 89,
      "unique_words": 67,
      "difficult_words": 23,
      "review_available": true
    },
    "recent_activity": [
      {
        "timestamp": "2025-11-21T14:15:00Z",
        "activity_type": "page_read",
        "page_number": 67,
        "details": "Completed section 2"
      },
      {
        "timestamp": "2025-11-21T14:10:00Z",
        "activity_type": "vocabulary_lookup",
        "word": "curtsied",
        "section": "p67_s2"
      }
    ]
  }
}
```

### HELP SYSTEM SCHEMAS

#### Help Request Creation
```json
// Request Schema
{
  "context": {
    "current_page": 89,
    "current_section": 2,
    "selected_text": "I do wish I hadn't cried so much!",
    "surrounding_text": "Everything is queer to-day. I do wish I hadn't cried so much!"
  },
  "help_category": "comprehension", // comprehension, vocabulary, analysis, motivation
  "difficulty_level": "confused",
  "user_notes": "I don't understand why Alice is crying or what's 'queer about today'",
  "urgency": "normal", // normal, urgent, critical
  "preferred_contact_method": "text", // text, voice_call, phone_consultation
  "availability": {
    "timezone": "America/New_York",
    "available_times": [
      {
        "date": "2025-11-21",
        "start_time": "15:00",
        "end_time": "17:00"
      }
    ]
  }
}

// Response Schema
{
  "success": true,
  "data": {
    "help_request": {
      "id": "help_req_789xyz123abc",
      "status": "pending", // pending, assigned, in_progress, resolved, closed
      "created_at": "2025-11-21T14:35:00Z",
      "estimated_response_time_minutes": 30,
      "assigned_consultant": null,
      "context_summary": {
        "book_title": "Alice's Adventures in Wonderland",
        "current_page": 89,
        "current_section": 2,
        "selected_text": "I do wish I hadn't cried so much!",
        "user_category": "comprehension"
      },
      "priority_score": 6.5, // 0-10 based on urgency, history, complexity
      "user_class_details": {
        "previous_help_requests": 3,
        "average_response_time_minutes": 25,
        "resolution_satisfaction_score": 8.5,
        "reading_progress_percentage": 45.2,
        "frequently_confused_concepts": ["symbolism", "emotional_motivation"]
      },
      "communication_channel": {
        "channel_type": "in_app_chat",
        "channel_id": "chat_456def789ghi",
        "status": "ready_for_consultant_connection"
      }
    },
    "next_steps": {
      "waiting_message": "A consultant will review your question shortly",
      "estimated_first_response": "2025-11-21T15:05:00Z",
      "availability_confirmation": true,
      "recommended_resources": [
        {
          "title": "Understanding Alice's Emotional Journey",
          "type": "reading_guide",
          "url": "/resources/alice-emotional-journey"
        }
      ]
    }
  }
}
```

#### Consultant Response Schema
```json
{
  "success": true,
  "data": {
    "consultant_info": {
      "id": "consultant_jane_smith_001",
      "name": "Jane Smith",
      "credentials": "M.Ed. Children's Literature",
      "specialties": ["19th_century_literature", "symbolism", "character_development"],
      "response_stats": {
        "total_responses": 342,
        "average_response_time_minutes": 12,
        "satisfaction_rating": 9.2,
        "resolution_rate_percentage": 97
      },
      "availability_status": "online",
      "timezone": "America/New_York"
    },
    "response": {
      "id": "response_123ghi789jkl",
      "help_request_id": "help_req_789xyz123abc",
      "timestamp": "2025-11-21T14:45:00Z",
      "response_type": "explanation", // explanation, question, resource, exercise
      "content": "Alice is experiencing the emotional weight of her transformation. The phrase 'queer about today' reflects her recognition that Wonderland operates by different rules than the normal world she's used to.",
      "simplification": "Alice is crying because everything is confusing and different to her.",
      "detailed_analysis": "This moment represents Alice's growing awareness that she's in a place where normal logic doesn't apply...",
      "additional_resources": [
        {
          "type": "text_excerpt",
          "title": "Coming-of-age in Wonderland",
          "content": "The theme of growing up and changing perspective...",
          "source": "literary_analysis_companion",
          "page_range": "45-47"
        }
      ],
      "follow_up_questions": [
        "How do you think you would feel if everything around you stopped making sense?",
        "What does Alice learn about herself through these confusing experiences?"
      ],
      "next_action": "Think about times when you've felt confused about your surroundings. How did you adapt?"
    },
    "conversation_context": {
      "session_duration_minutes": 15,
      "messages_exchanged": 4,
      "user_engagement_score": 8,
      "learning_outcomes_achieved": ["emotional_understanding", "character_motivation"],
      "recommended_session_extension": false,
      "resolution_probability": 0.85
    }
  }
}
```

## ERROR RESPONSE SCHEMAS

### Standard Error Structure
```json
{
  "success": false,
  "timestamp": "2025-11-21T14:30:00Z",
  "request_id": "req_failed_123def",
  "error": {
    "type": "validation_error",
    "code": "INVALID_PAGE_NUMBER",
    "message": "Page 250 doesn't exist in this book. Please use a page between 1 and 200.",
    "details": {
      "field": "page_number",
      "provided_value": 250,
      "valid_range": "1-200",
      "closest_valid_page": 200
    },
    "suggestions": [
      "Try page 200 for the final chapter",
      "Use the navigation buttons to browse available pages"
    ],
    "retry_allowed": true,
    "timestamp": "2025-11-21T14:30:00Z",
    "documentation_url": "/api/docs/errors#invalid_page_number",
    "fallback_options": {
      "type": "suggest_valid_range",
      "options": ["previous_page", "book_navigation", "random_page"]
    }
  },
  "metadata": {
    "error_timestamp": "2025-11-21T14:30:00Z",
    "failed_component": "page_service",
    "user_facing_severity": "low",
    "needs_consultant_attention": false
  }
}
```

### Database Error Schema
```json
{
  "success": false,
  "error": {
    "type": "database_error",
    "code": "CONNECTION_TIMEOUT",
    "message": "We're experiencing technical difficulties accessing the book content.",
    "technical_message": "SQLite connection timeout after 5000ms",
    "user_action": "Please try again in a few moments.",
    "retry_after_ms": 30000,
    "fallback_available": true,
    "fallback_type": "cached_content",
    "estimated_available_pages": 15,
    "contact_support": false
  },
  "metadata": {
    "error_category": "infrastructure",
    "recovery_time_estimate": "2-5 minutes",
    "affected_services": ["page_content", "glossary_lookup", "progress_tracking"]
  }
}
```

## METADATA SCHEMAS

### Reading Context Metadata
```json
{
  "reading_context": {
    "current_state": {
      "book_id": "book_alice_adventures_wonderland",
      "page_number": 67,
      "section_number": 2,
      "position_percentage": 33.5,
      "reading_progress_percentage": 33.5
    },
    "user_progression": {
      "pages_read_today": 8,
      "consecutive_reading_days": 5,
      "average_pages_per_session": 6.2,
      "average_session_duration_minutes": 22,
      "reading_velocity_wpm": 124
    },
    "engagement_metrics": {
      "help_requests_this_session": 2,
      "ai_interactions_this_session": 5,
      "vocabulary_lookups_this_session": 12,
      "bookmarks_this_session": 1,
      "session_start_time": "2025-11-21T13:45:00Z",
      "current_session_duration_minutes": 45
    },
    "learning_indicators": {
      "difficulty_identified": false,
      "conceptual_confusion_areas": ["symbolism"],
      "vocabulary_struggles": ["archaic_expressions", "literary_references"],
      "comprehension_level": "developing",
      "engagement_score": 7.8,
      "frustration_indicators": 0.2
    }
  }
}
```

### System Performance Metadata
```json
{
  "performance_metrics": {
    "response_time_ms": 234,
    "cache_stats": {
      "cache_hit": true,
      "cache_key": "page_67_alice_book",
      "cache_age_seconds": 1800
    },
    "database_stats": {
      "query_time_ms": 45,
      "records_returned": 8,
      "query_type": "SELECT",
      "index_usage": ["page_number_idx"]
    },
    "server_resources": {
      "memory_usage_mb": 128,
      "cpu_usage_percentage": 23,
      "concurrent_connections": 12
    }
  }
}
```

## VALIDATION RULES AND CONSTRAINTS

### Input Validation Schemas

#### Page Number Validation
```json
{
  "page_number": {
    "type": "integer",
    "minimum": 1,
    "maximum": 200,
    "error_message": "Page must be between 1 and 200"
  }
}
```

#### User Registration Validation
```json
{
  "email": {
    "type": "string",
    "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
    "max_length": 254,
    "error_message": "Please provide a valid email address"
  },
  "password": {
    "type": "string",
    "min_length": 8,
    "must_contain": ["number", "uppercase", "lowercase"],
    "error_message": "Password must be 8+ characters with letters and numbers"
  },
  "first_name": {
    "type": "string",
    "pattern": "^[a-zA-Z\\s]{2,50}$",
    "error_message": "First name should contain only letters"
  }
}
```

#### Book Access Code Validation
```json
{
  "access_code": {
    "type": "string",
    "pattern": "^[A-Z0-9]{6,8}$",
    "case_sensitive": false,
    "error_message": "Access code should be 6-8 characters (letters and numbers)"
  }
}
```

## DATA INTEGRITY SCHEMAS

### Book Content Structure Validation
```json
{
  "book_content_validation": {
    "page_word_count": {
      "min": 150,
      "max": 300,
      "average": 185,
      "standard_deviation": 25
    },
    "section_word_count": {
      "min": 35,
      "max": 50,
      "average": 42,
      "sections_per_page": {
        "min": 3,
        "max": 5,
        "average": 4.2
      }
    },
    "glossary_linking": {
      "sections_per_glossary_term": {
        "min": 1,
        "max": 15,
        "average": 3.4
      },
      "context_accuracy_minimum": 0.8,
      "linking_confidence_threshold": 0.85
    }
  }
}
```

This comprehensive schema specification ensures consistent, predictable API behavior while maintaining flexibility for future enhancements. All implementations must conform to these schemas for successful integration with the Alice Reader App ecosystem.
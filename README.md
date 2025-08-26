# Interactive Book Generation Pipeline

This document provides instructions on how to use the interactive, single-book generation pipeline.

## Overview

The pipeline is an interactive system that guides you through the creation of a single book. It prompts you for book details, generates an initial outline for your review, and then proceeds with a full 7-phase generation process to create the complete book and its marketing materials.

This new interactive workflow replaces the previous batch-processing system.

## How to Generate a New Book

Follow these steps to generate a new book.

### Step 1: Run the Interactive Pipeline

1.  **Open your terminal.**
2.  **Execute the main script**:
    ```bash
    bash run_sequential_pipeline.sh
    ```

### Step 2: Enter Book Details

The script will prompt you to enter the following details for your new book:

-   **Title**: The main title of the book.
-   **Memorable Phrase (Subtitle)**: A catchy subtitle or phrase.
-   **Category**: The genre or category (e.g., "AI & Business").
-   **Description**: A brief description of the book's content.

### Step 3: Duplicate Check

The system will check if a book with the same title already exists in the `output/books` directory. If a duplicate is found, it will ask if you want to overwrite it.

### Step 4: Review the Outline

After you provide the details, the system will generate a 1000-word outline and display a summary on the screen, including:

-   Title & Logline
-   Synopsis
-   A list of chapter titles

### Step 5: Confirm and Generate

The script will start a **2-minute countdown**. You have two options:

1.  **Do Nothing**: Let the timer run out, and the full 7-phase book generation will begin automatically.
2.  **Stop the Process**: Press `s` and then `Enter` to cancel the generation before the timer ends.

### Step 6: Access Your Generated Book

Once the process is complete, all generated files for your new book are stored in the `output/books/` directory inside a new `book_XXX` folder. This includes the outline, phase results, and the final book in multiple formats.